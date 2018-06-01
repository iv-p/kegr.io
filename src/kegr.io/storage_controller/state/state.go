package state

import (
	"io/ioutil"
	"log"

	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/model/keg"
	"kegr.io/storage_controller/model/state"
)

// StateService is responsible for comparing and keeping the state up to date
type StateService struct {
	IStateService

	kegByPath map[string]keg.IKeg
	kegByID   map[string]keg.IKeg
}

// IStateService is StateServices interface
type IStateService interface {
	GetKegByID(kegID string) (keg.IKeg, error)
	GetKegByPath(path string) (keg.IKeg, error)
	GetKegs() map[string]keg.IKeg

	// Keg operations
	CreateKeg(options keg.IOptions) (keg.IKeg, error)
	UpdateKeg(kegID string, options keg.IOptions) error
	DeleteKeg(kegID string) error

	// Liquid operations
	GetState() state.IState
	GetHash() ([]byte, error)
	Diff(kegs map[string]keg.IKeg) map[string]*keg.KegDiff
}

// NewStateService returns an initialised state service object
func NewStateService() IStateService {
	ss := &StateService{
		kegByPath: make(map[string]keg.IKeg),
		kegByID:   make(map[string]keg.IKeg),
	}

	files, err := ioutil.ReadDir(config.C.DataRoot)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			if keg, err := keg.FromDir(file.Name()); err == nil {
				ss.addKeg(keg)
			}
		}
	}

	return ss
}

// GetState returns the current state of the server
func (ss *StateService) GetState() state.IState {
	state := state.NewState()
	state.SetKegs(ss.kegByID)
	return state
}

// GetHash returns the hash of the state
func (ss *StateService) GetHash() ([]byte, error) {
	return ss.GetState().GetHash()
}
