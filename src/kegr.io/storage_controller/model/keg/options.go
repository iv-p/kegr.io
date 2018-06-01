package keg

import pbKeg "kegr.io/protobuf/model/storage/keg"

// Options hold the changeable data for a keg
type Options struct {
	name  string
	path  string
	cache int64
	gzip  bool
	IOptions
}

// IOptions is an interface
type IOptions interface {
	GetName() string
	SetName(name string)
	GetGzip() bool
	SetGzip(gzip bool)
	GetCache() int64
	SetCache(cache int64)
	GetPath() string
	SetPath(path string)
	Diff(other IOptions) IOptions

	ToProto() *pbKeg.Options
}

// NewOptions returns a new ILiquidOptions object
func NewOptions() IOptions {
	return &Options{
		cache: 0,
		gzip:  false,
	}
}

// OptionsFromProto converts the protobuf  options object to a model.options
func OptionsFromProto(lo *pbKeg.Options) IOptions {
	return &Options{
		name:  lo.Name,
		path:  lo.Path,
		cache: lo.Cache,
		gzip:  lo.Gzip,
	}
}

// Diff returns the difference with preference the other element
func (o *Options) Diff(other IOptions) IOptions {
	newOptions := &Options{}
	newOptions.SetName(o.GetName())
	newOptions.SetGzip(o.GetGzip())
	newOptions.SetCache(o.GetCache())
	newOptions.SetPath(o.GetPath())
	if o.GetName() != other.GetName() {
		newOptions.SetName(other.GetName())
	}
	if o.GetGzip() != other.GetGzip() {
		newOptions.SetGzip(other.GetGzip())
	}
	if o.GetCache() != other.GetCache() {
		newOptions.SetCache(other.GetCache())
	}
	if o.GetPath() != other.GetPath() {
		newOptions.SetPath(other.GetPath())
	}
	return newOptions
}

// GetName getter
func (o *Options) GetName() string {
	return o.name
}

// SetName setter
func (o *Options) SetName(name string) {
	o.name = name
}

// GetGzip getter
func (o *Options) GetGzip() bool {
	return o.gzip
}

// SetGzip setter
func (o *Options) SetGzip(gzip bool) {
	o.gzip = gzip
}

// GetCache getter
func (o *Options) GetCache() int64 {
	return o.cache
}

// SetCache setter
func (o *Options) SetCache(cache int64) {
	o.cache = cache
}

// GetPath getter
func (o *Options) GetPath() string {
	return o.path
}

// SetPath setter
func (o *Options) SetPath(path string) {
	o.path = path
}

// ToProto returns the proto representation of the object
func (o *Options) ToProto() *pbKeg.Options {
	return &pbKeg.Options{
		Name:  o.name,
		Path:  o.path,
		Cache: o.cache,
		Gzip:  o.gzip,
	}
}

func optionsFromProto(o *pbKeg.Options) IOptions {
	return &Options{
		name:  o.Name,
		path:  o.Path,
		cache: o.Cache,
		gzip:  o.Gzip,
	}
}
