package state

import (
	"crypto/md5"
	"sort"

	"github.com/golang/protobuf/proto"
	pbKeg "kegr.io/protobuf/model/storage/keg"
	pbState "kegr.io/protobuf/model/storage/state"
	"kegr.io/storage_controller/model/keg"
)

// State holds the current state of the server
type State struct {
	IState

	kegs map[string]keg.IKeg
}

// IState is the interface a State should implement
type IState interface {
	GetHash() ([]byte, error)
	ToProto() *pbState.State
	ToBytes() ([]byte, error)
	SetKegs(kegs map[string]keg.IKeg)
}

// NewState returns an initialised state object
func NewState() IState {
	return &State{
		kegs: make(map[string]keg.IKeg),
	}
}

// GetHash returns the combined hash of all kegs in the state
func (s *State) GetHash() ([]byte, error) {
	var keys []string
	for k := range s.kegs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	hash := md5.New()
	for _, id := range keys {
		bytes, err := s.kegs[id].GetStateHash()
		if err != nil {
			return nil, err
		}
		hash.Write(bytes)
	}
	return hash.Sum(nil), nil
}

// ToProto returns the proto in protobuf state
func (s *State) ToProto() *pbState.State {
	state := &pbState.State{
		Kegs: make(map[string]*pbKeg.Keg),
	}
	for id, keg := range s.kegs {
		state.Kegs[id] = keg.ToProto()
	}
	return state
}

// ToBytes returns the state in byte array form
func (s *State) ToBytes() ([]byte, error) {
	return proto.Marshal(s.ToProto())
}

// SetKegs sets the kegs
func (s *State) SetKegs(kegs map[string]keg.IKeg) {
	s.kegs = kegs
}
