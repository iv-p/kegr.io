package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"kegr.io/protobuf/model/storage/keg"

	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/merkle"
	"kegr.io/storage_controller/model"
)

// KegService holds information for the kegs currently
// found on the FS
type KegService struct {
	c         *config.Config
	kegByPath map[string]model.IKeg
	kegByID   map[string]model.IKeg
	IKegService
}

// IKegService is an interface
type IKegService interface {
	Init()
	CreateKeg(string, string) model.IKeg
	DeleteKeg(string) error

	GetLiquidByPath(string, string) (model.ILiquid, error)
	GetLiquidByIds(string, string) (model.ILiquid, error)
	AddLiquid(string, model.ILiquid, bool) error
	UpdateLiquid(string, string, model.ILiquidOptions) error
	DeleteLiquid(string, string) error
	GetLiquids(string) ([]model.ILiquidInfo, error)
	GetState() []byte
	GetProtoState() *keg.KegServiceState
	Compare(map[string]*model.Keg) map[string][]merkle.IContent
}

// NewKegService creates a new KegService object and
// tries to load existing kegs from the FS' dataRoot folder
func NewKegService(c *config.Config) *KegService {
	ks := &KegService{
		c:         c,
		kegByPath: make(map[string]model.IKeg),
		kegByID:   make(map[string]model.IKeg),
	}
	ks.Init()
	return ks
}

// Init will read the dataRoot folder and try to find existing
// keg folders to add to the kegService
func (ks *KegService) Init() {
	files, err := ioutil.ReadDir(ks.c.DataRoot)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if keg, err := model.LoadKegFromDir(file.Name(), ks.c); err == nil {
				ks.addKegToMaps(keg)
			}
		}
	}
}

func (ks *KegService) addKegToMaps(keg model.IKeg) {
	ks.kegByPath[keg.GetPath()] = keg
	ks.kegByID[keg.GetID()] = keg
}

func (ks *KegService) deleteKegFromMaps(keg model.IKeg) {
	delete(ks.kegByPath, keg.GetPath())
	delete(ks.kegByID, keg.GetID())
}

// CreateKeg craetes a new Keg and adds it to the list of existing kegs.
func (ks *KegService) CreateKeg(name string, path string) model.IKeg {
	keg := model.NewKeg("", name, path, ks.c)
	ks.addKegToMaps(keg)
	return keg
}

// DeleteKeg removes a Keg and deletes all the files in it
func (ks *KegService) DeleteKeg(kegID string) error {
	var keg model.IKeg
	var err error
	var exist bool

	if keg, exist = ks.kegByID[kegID]; !exist {
		return errors.New("Keg not found")
	}

	if err = os.RemoveAll(ks.getKegPath(keg)); err != nil {
		return errors.New("Could not delete keg folder")
	}

	ks.deleteKegFromMaps(keg)
	return nil
}

// GetLiquidByPath finds the corresponding keg by accessPath and then if found
// tries to load the corresponding liquid in that keg.
func (ks *KegService) GetLiquidByPath(kegAccessPath string, liquidAccessName string) (model.ILiquid, error) {
	if keg, kegExist := ks.kegByPath[kegAccessPath]; kegExist {
		return keg.GetLiquidByAccessName(liquidAccessName)
	}
	return nil, errors.New("Keg not found")
}

// GetLiquidByIds takes the keg and liquid ids and returns the liquid if found
func (ks *KegService) GetLiquidByIds(kegID string, liquidID string) (model.ILiquid, error) {
	if keg, kegExist := ks.kegByID[kegID]; kegExist {
		return keg.GetLiquidByID(liquidID)
	}
	return nil, errors.New("Keg not found")
}

// AddLiquid will look for the keg we want to add a liquid to and if found
// will create a new liquid file on the FS
func (ks *KegService) AddLiquid(kegID string, liquid model.ILiquid, force bool) error {
	if keg, exist := ks.kegByID[kegID]; exist {
		return keg.CreateLiquid(liquid, force)
	}
	return errors.New("Keg not found")
}

// UpdateLiquid gets an options object and updates the liquid's options
func (ks *KegService) UpdateLiquid(kegID string, liquidID string, options model.ILiquidOptions) error {
	if keg, exist := ks.kegByID[kegID]; exist {
		return keg.UpdateLiquid(liquidID, options)
	}
	return errors.New("Keg not found")
}

// DeleteLiquid finds the corrent keg to delete the liquid from and if found tries to
// delete the liquid with that id.
func (ks *KegService) DeleteLiquid(kegID string, liquidID string) error {
	if keg, exist := ks.kegByID[kegID]; exist {
		return keg.DeleteLiquid(liquidID)
	}
	return errors.New("Keg not found")
}

// GetLiquids returns the liquids that are in this keg's info
func (ks *KegService) GetLiquids(kegID string) ([]model.ILiquidInfo, error) {
	if keg, exist := ks.kegByID[kegID]; exist {
		return keg.GetAllLiquidInfo(), nil
	}
	return nil, errors.New("Keg not found")
}

// GetState returns the hash of the kegs
func (ks *KegService) GetState() []byte {
	var keys []string
	for k := range ks.kegByID {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	hash := md5.New()
	for _, id := range keys {
		hash.Write(ks.kegByID[id].GetHash())
	}
	return hash.Sum(nil)
}

// GetProtoState returns the proto representation of a map with all of
// the merkle trees (one per keg)
func (ks *KegService) GetProtoState() *keg.KegServiceState {
	res := make(map[string]*keg.PbKeg)
	for key, keg := range ks.kegByID {
		res[key] = keg.GetState()
	}

	return &keg.KegServiceState{
		Kegs: res,
	}
}

// Compare returns the timestamp when we last
// had a bodification on any bucket
func (ks *KegService) Compare(kegs map[string]*model.Keg) map[string][]merkle.IContent {
	diff := make(map[string][]merkle.IContent)

	for id, keg := range kegs {
		localKeg, exist := ks.kegByID[id]
		if !exist {
			localKeg = model.NewKeg(keg.GetID(), keg.GetName(), keg.GetPath(), ks.c)
			ks.addKegToMaps(localKeg)
		}

		d := localKeg.Compare(keg.GetTree())
		diff[id] = d
	}

	return diff
}

func (ks *KegService) getKegPath(keg model.IKeg) string {
	return fmt.Sprintf("%v/%v", ks.c.DataRoot, keg.GetID())
}
