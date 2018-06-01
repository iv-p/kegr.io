package keg

import (
	"errors"

	"kegr.io/storage_controller/model/liquid"
)

// AddLiquid inserts a liquid in all data structures and makes it available in this keg
func (k *Keg) AddLiquid(info liquid.IInfo) error {
	return k.updateLiquidInfo(info)
}

// UpdateLiquid updates the liquids options and updates the merkle tree
func (k *Keg) UpdateLiquid(info liquid.IInfo) error {
	return k.updateLiquidInfo(info)
}

// DeleteLiquid receives a liquidID and checks if this liquid is in the keg,
// after which it marks it as deleted
func (k *Keg) DeleteLiquid(liquidID string) error {
	info, exist := k.liquidInfo[liquidID]
	if !exist || info.IsDeleted() {
		return errors.New("File not found")
	}

	info.SetDeleted(true)

	return k.updateLiquidInfo(info)
}

// GetLiquidIDByAccessName does a lookup in the kegs LiquidByAccessName map to
// find the liquidID of the corresponding file and loads it form disk.
func (k *Keg) GetLiquidIDByAccessName(liquidAccessName string) (string, error) {
	liquidID, exist := k.liquidByAccessName[liquidAccessName]
	if !exist {
		return "", errors.New("File not found")
	}
	return liquidID, nil
}

// GetLiquidInfoByID returns the liquid info for that id
func (k *Keg) GetLiquidInfoByID(liquidID string) (liquid.IInfo, error) {
	info, exist := k.liquidInfo[liquidID]
	if !exist {
		return nil, errors.New("File not found")
	}
	return info, nil
}

// GetLiquids returns all the liquids in this keg
func (k *Keg) GetLiquids() map[string]liquid.IInfo {
	return k.liquidInfo
}

func (k *Keg) updateLiquidInfo(info liquid.IInfo) error {
	oldInfo, exist := k.liquidInfo[info.GetID()]
	if exist {
		delete(k.liquidByAccessName, oldInfo.GetAccessName())
		delete(k.liquidInfo, oldInfo.GetID())
	}

	// Add new entries
	k.liquidInfo[info.GetID()] = info
	k.liquidByAccessName[info.GetAccessName()] = info.GetID()
	merkleTreeLiquid, err := liquid.NewMerkleTreeLiquid(info)
	if err != nil {
		return err
	}
	return k.merkleTree.Add(merkleTreeLiquid)
}
