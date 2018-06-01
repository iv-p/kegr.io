package liquid

import pbLiquid "kegr.io/protobuf/model/storage/liquid"

// Options options holds all the changealbe options of a liquid
type Options struct {
	name  string
	ext   string
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
	GetExt() string
	SetExt(ext string)
}

// NewOptions returns a new ILiquidOptions object
func NewOptions() IOptions {
	return &Options{}
}

func OptionsFromProto(lo *pbLiquid.Options) IOptions {
	return &Options{
		name:  lo.Name,
		ext:   lo.Ext,
		cache: lo.Cache,
		gzip:  lo.Gzip,
	}
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

// GetExt getter
func (o *Options) GetExt() string {
	return o.ext
}

// SetExt setter
func (o *Options) SetExt(ext string) {
	o.ext = ext
}
