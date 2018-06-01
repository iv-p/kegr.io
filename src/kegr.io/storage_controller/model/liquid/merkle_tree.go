package liquid

import (
	"crypto/md5"

	pbMerkle "kegr.io/protobuf/model/merkle"
	"kegr.io/storage_controller/merkle"
)

// MerkleTreeLiquid is the representation of a liquid in the merkle tree
type MerkleTreeLiquid struct {
	merkle.IContent
	IMerkleTreeLiquid

	id          []byte
	hash        []byte
	lastUpdated int64
	deleted     bool
}

// IMerkleTreeLiquid is an interface
type IMerkleTreeLiquid interface {
	NewMerkleTreeLiquid(liquid ILiquid) (IMerkleTreeLiquid, error)

	GetID() []byte
	SetID(id []byte)
	GetHash() []byte
	SetHash(hash []byte)
	GetLastUpdated() int64
	SetLastUpdated(lastUpdated int64)
	GetProto() *pbMerkle.Content

	IsDeleted() bool
	SetDeleted(deleted bool)
}

// NewEmptyMerkleTreeLiquid returns an empty merkle tree liquid object
func NewEmptyMerkleTreeLiquid() *MerkleTreeLiquid {
	return &MerkleTreeLiquid{}
}

// NewMerkleTreeLiquid copies necessary infromation from a liquid to
// create its representation in the merkle tree
func NewMerkleTreeLiquid(info IInfo) (IMerkleTreeLiquid, error) {
	var bytes []byte
	var err error

	if bytes, err = info.ToBytes(); err != nil {
		return nil, err
	}

	hash := md5.New()
	hash.Write(bytes)

	return &MerkleTreeLiquid{
		id:          []byte(info.GetID()),
		hash:        hash.Sum(nil),
		lastUpdated: info.GetLastUpdated(),
		deleted:     info.IsDeleted(),
	}, nil
}

// GetProto returns the proto representation of the liquid
func (mtl *MerkleTreeLiquid) GetProto() *pbMerkle.Content {
	return &pbMerkle.Content{
		ID:          []byte(mtl.id),
		Hash:        mtl.hash,
		LastUpdated: mtl.lastUpdated,
		Deleted:     mtl.deleted,
	}
}

// GetID getter
func (mtl *MerkleTreeLiquid) GetID() []byte {
	return []byte(mtl.id)
}

// SetID setter
func (mtl *MerkleTreeLiquid) SetID(id []byte) {
	mtl.id = id
}

// GetHash getter
func (mtl *MerkleTreeLiquid) GetHash() []byte {
	return mtl.hash
}

// SetHash setter
func (mtl *MerkleTreeLiquid) SetHash(hash []byte) {
	mtl.hash = hash
}

// GetLastUpdated getter
func (mtl *MerkleTreeLiquid) GetLastUpdated() int64 {
	return mtl.lastUpdated
}

// SetLastUpdated setter
func (mtl *MerkleTreeLiquid) SetLastUpdated(lastUpdated int64) {
	mtl.lastUpdated = lastUpdated
}

// IsDeleted getter
func (mtl *MerkleTreeLiquid) IsDeleted() bool {
	return mtl.deleted
}

// SetDeleted setter
func (mtl *MerkleTreeLiquid) SetDeleted(deleted bool) {
	mtl.deleted = deleted
}
