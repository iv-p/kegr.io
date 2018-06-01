package keg

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	pbKeg "kegr.io/protobuf/model/storage/keg"

	"github.com/golang/protobuf/proto"
	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/merkle"
	"kegr.io/storage_controller/model/liquid"
)

const (
	kegFile         = ".keg"
	liquidExtension = "liquid"
	merkleTreeDepth = 16
)

// Keg holds information about a keg including quick access maps for liquidAccessName to liquidID
type Keg struct {
	IKeg

	id          string
	options     IOptions
	deleted     bool
	lastUpdated int64

	liquidByAccessName map[string]string
	liquidInfo         map[string]liquid.IInfo
	merkleTree         merkle.ITree
}

// IKeg is an interface
type IKeg interface {
	AddLiquid(info liquid.IInfo) error
	UpdateLiquid(info liquid.IInfo) error
	DeleteLiquid(liquidID string) error
	GetLiquidIDByAccessName(liquidAccessName string) (string, error)
	GetLiquidInfoByID(liquidID string) (liquid.IInfo, error)
	GetLiquids() map[string]liquid.IInfo

	ToBytes() ([]byte, error)
	ToProto() *pbKeg.Keg
	ToDir() error

	GetStateHash() ([]byte, error)
	Diff(other IKeg) *KegDiff

	// Getters and setters
	GetID() string
	GetOptions() IOptions
	SetOptions(options IOptions)
	GetTree() merkle.ITree
	GetLastUpdated() int64
	IsDeleted() bool
	SetDeleted(deleted bool)
	GetInfo() *Info
}

// ToBytes returns the bytes array representation of the object
func (k *Keg) ToBytes() ([]byte, error) {
	return proto.Marshal(k.toKegFile())
}

// ToProto returns the protobuf representation of this keg
func (k *Keg) ToProto() *pbKeg.Keg {
	return &pbKeg.Keg{
		Id:          k.id,
		Options:     k.options.ToProto(),
		Tree:        k.merkleTree.ToProto(),
		Deleted:     k.deleted,
		LastUpdated: k.lastUpdated,
	}
}

// ToDir saves the keg contents to the FS
func (k *Keg) ToDir() error {
	var content []byte
	var err error

	dir := fmt.Sprintf("%s/%s", config.C.DataRoot, k.id)
	kegFile := fmt.Sprintf("%s/%s", dir, config.C.KegFile)

	if err = os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if content, err = k.ToBytes(); err != nil {
		return err
	}

	return ioutil.WriteFile(kegFile, content, 0644)
}

// GetStateHash returns the bytes array representation of the object
func (k *Keg) GetStateHash() ([]byte, error) {
	bytes, err := k.ToBytes()
	if err != nil {
		return nil, err
	}

	hash := md5.New()
	hash.Write(bytes)
	hash.Write(k.merkleTree.Hash())
	return hash.Sum(nil), nil
}

// Diff returns the difference between this and another keg
func (k *Keg) Diff(other IKeg) *KegDiff {
	kd := NewKegDiff()

	// Compare options
	if other.GetLastUpdated() > k.GetLastUpdated() {
		kd.Options = other.GetOptions()
	}

	// Compare content
	kd.Content = k.GetTree().Diff(other.GetTree())

	return kd
}

func (k *Keg) GetInfo() *Info {
	return &Info{
		ID:          k.id,
		Deleted:     k.deleted,
		LastUpdated: k.lastUpdated,
		Name:        k.options.GetName(),
		Path:        k.options.GetPath(),
		Cache:       k.options.GetCache(),
		Gzip:        k.options.GetGzip(),
	}
}

func (k *Keg) toKegFile() *pbKeg.KegFile {
	return &pbKeg.KegFile{
		Id:          k.id,
		Options:     k.options.ToProto(),
		Deleted:     k.deleted,
		LastUpdated: k.lastUpdated,
	}
}
