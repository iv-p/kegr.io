package keg

import (
	"errors"
	"fmt"
	"time"

	pbKeg "kegr.io/protobuf/model/storage/keg"

	"github.com/golang/protobuf/proto"
	"kegr.io/storage_controller/merkle"
	"kegr.io/storage_controller/model/liquid"
	"kegr.io/storage_controller/util"
)

const (
	kegFile         = ".keg"
	liquidExtension = "liquid"
	merkleTreeDepth = 16
)

// Keg holds information about a keg including quick access maps for liquidAccessName to liquidID
type Keg struct {
	IKeg

	id      string
	options IOptions

	liquidByAccessName map[string]string
	liquidInfo         map[string]liquid.IInfo
	merkleTree         merkle.ITree
}

// IKeg is an interface
type IKeg interface {
	AddLiquid(liquid.IInfo)
	UpdateLiquid(liquid.IInfo)
	DeleteLiquid(liquidID string)

	GetLiquidIDByAccessName(liquidAccessName string)
	GetLiquidInfoByID(liquidID string)

	GetID() string
	GetName() string
	GetPath() string
	GetHash() []byte
	GetState() *pbKeg.Keg
}

// NewKeg initialises a new Keg object and creates the necessary dir and .keg file
// on the FS
func NewKeg(id string, options IOptions) IKeg {
	if len(id) == 0 {
		id = util.ID()
	}
	return &Keg{
		id:                 id,
		options:            options,
		liquidByAccessName: make(map[string]string),
		liquidInfo:         make(map[string]liquid.IInfo),
		merkleTree:         merkle.NewTree(merkleTreeDepth),
	}
}

// CreateLiquid receives a liquid without an ID (from the upload form).
// It assigns an ID and saves the liquid to file while modifying the access map for
// O(1) access by accessName.
func (k *Keg) CreateLiquid(liq liquid.ILiquid, force bool) error {
	var l liquid.IInfo
	if !force {
		if lID, exist := k.liquidByAccessName[liq.GetAccessName()]; exist {
			if l, exist = k.liquidInfo[lID]; exist && !l.IsDeleted() {
				return errors.New("File already exists")
			}
		}
	}

	if len(liq.GetID()) == 0 {
		id := l.GetID()
		if !l.IsDeleted() {
			id = util.ID()
		}
		liq.SetID(id)
	}

	return k.updateStateAndSaveLiquid(liq, true, !force)
}

// UpdateLiquid updates the liquids options and updates the merkle tree
func (k *Keg) UpdateLiquid(id string, options liquid.IOptions) error {
	var liquid liquid.ILiquid
	var err error

	if liquid, err = LiquidFromFile(k.getLiquidPath(id)); err != nil {
		return errors.New("File not found")
	}

	liquid.SetOptions(options)
	return k.updateStateAndSaveLiquid(liquid, true, true)
}

// DeleteLiquid receives a liquidID and checks if this liquid is in the keg,
// after which it deletes it if it is
func (k *Keg) DeleteLiquid(liquidID string) error {
	if l, exist := k.liquidInfo[liquidID]; !exist || l.IsDeleted() {
		return errors.New("File not found")
	}

	liquidFile := k.getLiquidPath(liquidID)

	var liquid liquid.ILiquid
	var err error

	for liquid, err = LiquidFromFile(liquidFile); err != nil; {
		return errors.New("File not found")
	}

	liquid.SetContent(nil)
	liquid.SetDeleted(true)

	return k.updateStateAndSaveLiquid(liquid, true, true)
}

func (k *Keg) updateStateAndSaveLiquid(liquid liquid.ILiquid, saveToFile bool, updateTimestamp bool) error {
	if updateTimestamp {
		liquid.SetLastUpdated(time.Now().Unix())
	}

	if saveToFile {
		err := liquid.ToFile(k.getLiquidPath(liquid.GetID()))
		if err != nil {
			return err
		}
	}

	k.liquidInfo[liquid.GetID()] = liquid.GetLiquidInfo()
	k.liquidByAccessName[liquid.GetAccessName()] = liquid.GetID()

	merkleTreeLiquid, err := liquid.NewMerkleTreeLiquid(liquid)
	if err != nil {
		return err
	}

	return k.merkleTree.Add(merkleTreeLiquid)
}

// GetLiquidByAccessName does a lookup in the kegs LiquidByAccessName map to
// find the liquidID of the corresponding file and loads it form disk.
func (k *Keg) GetLiquidByAccessName(accessName string) (liquid.ILiquid, error) {
	liquidID, exists := k.liquidByAccessName[accessName]
	if exists {
		if liquidOptions, exist := k.liquidInfo[liquidID]; exist && !liquidOptions.IsDeleted() {
			return LiquidFromFile(k.getLiquidPath(liquidID))
		}
	}
	return nil, errors.New("File not found")
}

// GetLiquidByID loads the liduid from disk and returns it.
func (k *Keg) GetLiquidByID(liquidID string) (liquid.ILiquid, error) {
	return LiquidFromFile(k.getLiquidPath(liquidID))
}

// GetLiquidInfo does a lookup based on ID for a liquid info and returns an
// error if not found
func (k *Keg) GetLiquidInfo(liquidID string) (liquid.IInfo, error) {
	if info, exist := k.liquidInfo[liquidID]; exist {
		return info, nil
	}
	return LiquidInfo{}, errors.New("File not found")
}

// GetAllLiquidInfo returns all liquid infos for the files in this keg
func (k *Keg) GetAllLiquidInfo() []liquid.IInfo {
	var infos []ILiquidInfo
	for _, v := range k.liquidInfo {
		infos = append(infos, v)
	}
	return infos
}

func (k *Keg) fromBytes(b []byte) error {
	h := &keg.PbKeg{}
	if err := proto.Unmarshal(b, h); err != nil {
		return err
	}

	k.id = h.GetId()
	k.path = h.GetPath()
	k.name = h.GetName()
	k.liquidByAccessName = make(map[string]string)

	return nil
}

// Compare compares this keg's merkle tree with another merkle tree and returns the diff
func (k *Keg) Compare(other merkle.ITree) []merkle.IContent {
	diff, _ := k.merkleTree.Diff(other)
	return diff
}

// GetID getter
func (k *Keg) GetID() string {
	return k.id
}

// GetName getter
func (k *Keg) GetName() string {
	return k.name
}

// GetPath getter
func (k *Keg) GetPath() string {
	return k.path
}

// GetTree getter
func (k *Keg) GetTree() merkle.ITree {
	return k.merkleTree
}

// GetHash getter
func (k *Keg) GetHash() []byte {
	return k.merkleTree.Hash()
}

// GetState returns the underlying merkle tree
func (k *Keg) GetState() *pbKeg.Keg {
	return &pbKeg.Keg{
		Id:   k.id,
		Name: k.name,
		Path: k.path,
		Tree: k.merkleTree.ToProto(),
	}
}

// KegFromProto converts a proto keg to a model.Keg
func KegFromProto(k *pbKeg.Keg) *Keg {
	return &Keg{
		id:         k.Id,
		name:       k.Name,
		path:       k.Path,
		merkleTree: merkle.TreeFromProto(k.Tree, liquid.NewEmptyMerkleTreeLiquid),
	}
}

func (k *Keg) toBytes() ([]byte, error) {
	return proto.Marshal(&keg.PbKeg{
		Id:   k.id,
		Path: k.path,
		Name: k.name,
	})
}

func (k *Keg) getLiquidPath(liquidID string) string {
	return fmt.Sprintf("%v/%v/%v.%v", k.c.DataRoot, k.id, liquidID, liquidExtension)
}
