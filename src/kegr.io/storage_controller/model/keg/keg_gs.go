package keg

import (
	"fmt"
	"time"

	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/merkle"
	"kegr.io/storage_controller/model/liquid"
)

// GetID getter
func (k *Keg) GetID() string {
	return k.id
}

// GetOptions getter
func (k *Keg) GetOptions() IOptions {
	return k.options
}

// GetLastUpdated getter
func (k *Keg) GetLastUpdated() int64 {
	return k.lastUpdated
}

// SetOptions setter
func (k *Keg) SetOptions(options IOptions) {
	k.lastUpdated = time.Now().Unix()
	k.options = options
	k.ToDir()
}

// GetTree getter
func (k *Keg) GetTree() merkle.ITree {
	return k.merkleTree
}

// IsDeleted getter
func (k *Keg) IsDeleted() bool {
	return k.deleted
}

// SetDeleted setter
func (k *Keg) SetDeleted(deleted bool) {
	k.lastUpdated = time.Now().Unix()
	k.deleted = deleted
	for _, li := range k.GetLiquids() {
		liquidFile := fmt.Sprintf("%s/%s/%s.%s", config.C.DataRoot, k.id, li.GetID(), config.C.LiquidExtension)
		if liquid, err := liquid.FromFile(liquidFile); err == nil {
			liquid.SetDeleted(true)
			k.updateLiquidInfo(liquid.GetLiquidInfo())
			liquid.ToFile(fmt.Sprintf("%s/%s", config.C.DataRoot, k.id))
		}
	}
	k.ToDir()
}
