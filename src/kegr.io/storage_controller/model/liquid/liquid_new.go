package liquid

import (
	"errors"
	"io/ioutil"
	"log"
	"strings"

	"github.com/golang/protobuf/proto"
	pbLiquid "kegr.io/protobuf/model/storage/liquid"
	"kegr.io/storage_controller/config"
)

// NewLiquid returns a new ILiquid object
func NewLiquid() ILiquid {
	return &Liquid{}
}

// FromBytes returns the model.Liquid representation of a byte array
func FromBytes(bytes []byte) (ILiquid, error) {
	liq := &pbLiquid.Liquid{}
	if err := proto.Unmarshal(bytes, liq); err != nil {
		return nil, err
	}

	return FromProto(liq), nil
}

// FromProto returns the model.Liquid representation of a proto liquid
func FromProto(proto *pbLiquid.Liquid) ILiquid {
	return &Liquid{
		id:          proto.ID,
		fileHash:    proto.FileHash,
		content:     proto.Content,
		size:        proto.Size,
		lastUpdated: proto.LastUpdated,
		deleted:     proto.Deleted,
		options:     OptionsFromProto(proto.Options),
	}
}

// FromFile loads a liquid model from a specified file.
func FromFile(file string) (ILiquid, error) {
	if !strings.HasSuffix(file, config.C.LiquidExtension) {
		return nil, errors.New("Invalid file extension")
	}

	var content []byte
	var err error

	if content, err = ioutil.ReadFile(file); err != nil {
		return nil, err
	}

	liq := &pbLiquid.Liquid{}

	if err := proto.Unmarshal(content, liq); err != nil {
		return nil, err
	}

	log.Printf("Loaded file %v\n", file)

	return FromProto(liq), nil
}
