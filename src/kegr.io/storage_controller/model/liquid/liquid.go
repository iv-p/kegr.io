package liquid

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	pbLiquid "kegr.io/protobuf/model/storage/liquid"
	"kegr.io/storage_controller/config"
)

// Liquid holds a Liquid's information
type Liquid struct {
	id          string
	fileHash    []byte
	content     []byte
	size        int64
	lastUpdated int64
	deleted     bool
	options     IOptions
	ILiquid
}

// ILiquid is an interface
type ILiquid interface {
	ToBytes() ([]byte, error)
	ToProto() *pbLiquid.Liquid
	ToFile(file string) error

	GetAccessName() string
	GetLiquidInfo() IInfo

	// Getters and setters
	GetID() string
	SetID(id string)
	GetFileHash() []byte
	SetFileHash(fileHash []byte)
	GetContent() []byte
	SetContent(content []byte)
	GetSize() int64
	SetSize(size int64)
	GetLastUpdated() int64
	SetLastUpdated(lastUpdated int64)
	IsDeleted() bool
	SetDeleted(deleted bool)
	GetOptions() IOptions
	SetOptions(options IOptions)
}

// ToBytes serializes the liquid in bytes array
func (l *Liquid) ToBytes() ([]byte, error) {
	return proto.Marshal(l.ToProto())
}

// ToProto returns the proto version of a liquid
func (l *Liquid) ToProto() *pbLiquid.Liquid {
	return &pbLiquid.Liquid{
		ID:          l.id,
		FileHash:    l.fileHash,
		Content:     l.content,
		Size:        l.size,
		LastUpdated: l.lastUpdated,
		Deleted:     l.deleted,
		Options: &pbLiquid.Options{
			Name:  l.options.GetName(),
			Ext:   l.options.GetExt(),
			Cache: l.options.GetCache(),
			Gzip:  l.options.GetGzip(),
		},
	}
}

// ToFile writes the liquid file to the fs in the specified path.
func (l *Liquid) ToFile(path string) error {
	var content []byte
	var err error

	if content, err = l.ToBytes(); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s.%s", path, l.id, config.C.LiquidExtension), content, 0644); err != nil {
		return err
	}

	return nil
}

// GetAccessName returns the string file name that one can
// use to access this resource via the cdn link
func (l *Liquid) GetAccessName() string {
	return fmt.Sprintf("%v.%x.%v", l.options.GetName(), l.fileHash, l.options.GetExt())
}

// GetLiquidInfo converts the information about a liquid in a
// liquidinfo object to store in memory
func (l *Liquid) GetLiquidInfo() IInfo {
	return &Info{
		ID:          l.id,
		FileHash:    l.fileHash,
		Size:        l.size,
		Name:        l.options.GetName(),
		Ext:         l.options.GetExt(),
		Cache:       l.options.GetCache(),
		Gzip:        l.options.GetGzip(),
		Deleted:     l.deleted,
		LastUpdated: l.lastUpdated,
	}
}
