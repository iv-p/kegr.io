package state

import (
	"errors"
	"fmt"

	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/model/keg"
)

// CreateKeg craetes a new Keg and adds it to the list of existing kegs.
func (ss *StateService) CreateKeg(options keg.IOptions) (keg.IKeg, error) {
	keg := keg.NewKeg(options)
	err := ss.addKeg(keg)
	if err != nil {
		return nil, err
	}
	keg.ToDir()
	return keg, nil
}

// UpdateKeg updates the options of a keg that is tracked
func (ss *StateService) UpdateKeg(kegID string, options keg.IOptions) error {
	keg, exist := ss.kegByID[kegID]
	if !exist {
		return errors.New("Keg not found")
	}

	delete(ss.kegByPath, keg.GetOptions().GetPath())
	keg.SetOptions(options)
	ss.kegByPath[keg.GetOptions().GetPath()] = keg

	return nil
}

// DeleteKeg removes the underlying FS folder
func (ss *StateService) DeleteKeg(kegID string) error {
	keg, exist := ss.kegByID[kegID]
	if !exist {
		return errors.New("Keg not found")
	}

	keg.SetDeleted(true)

	return nil
}

// GetKegByID returns the corresponding keg with that id or an error if not found
func (ss *StateService) GetKegByID(kegID string) (keg.IKeg, error) {
	keg, exist := ss.kegByID[kegID]
	if !exist {
		return nil, errors.New("Keg not found")
	}
	return keg, nil
}

// GetKegByPath returns the corresponding keg with that path or an error if not found
func (ss *StateService) GetKegByPath(path string) (keg.IKeg, error) {
	keg, exist := ss.kegByPath[path]
	if !exist {
		return nil, errors.New("Keg not found")
	}
	return keg, nil
}

// GetKegs returns all kegs in a map where the key is the kegID
func (ss *StateService) GetKegs() map[string]keg.IKeg {
	return ss.kegByID
}

func (ss *StateService) getKegPath(id string) string {
	return fmt.Sprintf("%v/%v", config.C.DataRoot, id)
}

func (ss *StateService) addKeg(keg keg.IKeg) error {
	if _, exist := ss.kegByPath[keg.GetOptions().GetPath()]; exist {
		return errors.New("keg already exist")
	}
	ss.kegByPath[keg.GetOptions().GetPath()] = keg
	ss.kegByID[keg.GetID()] = keg
	return nil
}
