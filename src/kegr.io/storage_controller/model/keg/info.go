package keg

import "kegr.io/protobuf/model/storage/keg"

// Info holds information for a keg
type Info struct {
	ID          string
	Deleted     bool
	LastUpdated int64
	Name        string
	Path        string
	Cache       int64
	Gzip        bool
}

func (i *Info) ToProto() *keg.Info {
	return &keg.Info{
		Id:          i.ID,
		Deleted:     i.Deleted,
		LastUpdated: i.LastUpdated,
		Name:        i.Name,
		Path:        i.Path,
		Cache:       i.Cache,
		Gzip:        i.Gzip,
	}
}
