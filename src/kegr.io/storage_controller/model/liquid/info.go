package liquid

import (
	"gopkg.in/mgo.v2/bson"
	"kegr.io/protobuf/model/storage/liquid"
)

// Info holds information of a liquid that is going to be passed to
// the front end. It is kept in memory for fast operation.
type Info struct {
	IInfo

	ID          string
	FileHash    []byte
	Size        int64
	Name        string
	Ext         string
	Cache       int64
	Gzip        bool
	Deleted     bool
	AccessName  string
	LastUpdated int64
}

// IInfo is an interface
type IInfo interface {
	GetID() string
	SetID(id string)
	GetFileHash() []byte
	SetFileHash(fileHash []byte)
	GetSize() int64
	SetSize(size int64)
	GetName() string
	SetName(name string)
	GetExt() string
	SetExt(ext string)
	GetCache() int64
	SetCache(cache int64)
	IsGzip() bool
	SetGzip(gzip bool)
	IsDeleted() bool
	SetDeleted(deleted bool)
	GetAccessName() string
	SetAccessName(accessName string)
	GetLastUpdated() int64
	SetLastUpdated(lastUpdated int64)

	ToBytes() ([]byte, error)
	ToProto() *liquid.Info
}

// GetID getter
func (i *Info) GetID() string {
	return i.ID
}

// SetID setter
func (i *Info) SetID(id string) {
	i.ID = id
}

// GetFileHash getter
func (i *Info) GetFileHash() []byte {
	return i.FileHash
}

// SetFileHash setter
func (i *Info) SetFileHash(fileHash []byte) {
	i.FileHash = fileHash
}

// GetSize getter
func (i *Info) GetSize() int64 {
	return i.Size
}

// SetSize setter
func (i *Info) SetSize(size int64) {
	i.Size = size
}

// GetName getter
func (i *Info) GetName() string {
	return i.Name
}

// SetName setter
func (i *Info) SetName(name string) {
	i.Name = name
}

// GetExt getter
func (i *Info) GetExt() string {
	return i.Ext
}

// SetExt setter
func (i *Info) SetExt(ext string) {
	i.Ext = ext
}

// GetCache getter
func (i *Info) GetCache() int64 {
	return i.Cache
}

// SetCache setter
func (i *Info) SetCache(cache int64) {
	i.Cache = cache
}

// IsGzip getter
func (i *Info) IsGzip() bool {
	return i.Gzip
}

// SetGzip setter
func (i *Info) SetGzip(gzip bool) {
	i.Gzip = gzip
}

// IsDeleted getter
func (i *Info) IsDeleted() bool {
	return i.Deleted
}

// SetDeleted setter
func (i *Info) SetDeleted(deleted bool) {
	i.Deleted = deleted
}

// GetAccessName getter
func (i *Info) GetAccessName() string {
	return i.AccessName
}

// SetAccessName setter
func (i *Info) SetAccessName(accessName string) {
	i.AccessName = accessName
}

// GetLastUpdated getter
func (i *Info) GetLastUpdated() int64 {
	return i.LastUpdated
}

// SetLastUpdated setter
func (i *Info) SetLastUpdated(lastUpdated int64) {
	i.LastUpdated = lastUpdated
}

// ToBytes returns the byte array representation of the object
func (i *Info) ToBytes() ([]byte, error) {
	return bson.Marshal(*i)
}

// ToProto returns the proto info object
func (i *Info) ToProto() *liquid.Info {
	return &liquid.Info{
		Id:          i.ID,
		FileHash:    i.FileHash,
		Size:        i.Size,
		Name:        i.Name,
		Ext:         i.Ext,
		Cache:       i.Cache,
		Gzip:        i.Gzip,
		Deleted:     i.Deleted,
		AccessName:  i.AccessName,
		LastUpdated: i.LastUpdated,
	}
}
