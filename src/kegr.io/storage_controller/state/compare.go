package state

import (
	"kegr.io/storage_controller/model/keg"
)

// Diff return the state difference with another state
func (ss *StateService) Diff(kegs map[string]keg.IKeg) map[string]*keg.KegDiff {
	diff := make(map[string]*keg.KegDiff)

	for id, k := range kegs {
		localKeg, exist := ss.kegByID[id]
		if !exist {
			localKeg = keg.NewKegWithID(k.GetID(), k.GetOptions())
			ss.addKeg(localKeg)
			localKeg.ToDir()
		}

		d := localKeg.Diff(k)
		diff[id] = d
	}

	return diff
}
