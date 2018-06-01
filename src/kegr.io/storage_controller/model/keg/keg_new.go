package keg

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	pbKeg "kegr.io/protobuf/model/storage/keg"
	"kegr.io/storage_controller/config"
	"kegr.io/storage_controller/merkle"
	"kegr.io/storage_controller/model/liquid"
	"kegr.io/storage_controller/util"
)

// NewKegWithID initialises a new Keg object with the specified id
func NewKegWithID(id string, options IOptions) IKeg {
	return &Keg{
		id:                 id,
		options:            options,
		liquidByAccessName: make(map[string]string),
		liquidInfo:         make(map[string]liquid.IInfo),
		merkleTree:         merkle.NewTree(merkleTreeDepth),
		lastUpdated:        time.Now().Unix(),
		deleted:            false,
	}
}

// NewKeg initialises a new Keg object and assigns it a random id
func NewKeg(options IOptions) IKeg {
	return NewKegWithID(util.ID(), options)
}

// NewEmptyKeg initialises a new Keg object and creates the necessary dir and .keg file
// on the FS
func NewEmptyKeg() IKeg {
	return &Keg{
		liquidByAccessName: make(map[string]string),
		liquidInfo:         make(map[string]liquid.IInfo),
		merkleTree:         merkle.NewTree(merkleTreeDepth),
		lastUpdated:        time.Now().Unix(),
		deleted:            false,
	}
}

// FromBytes converts a byte array to a model.Keg
func FromBytes(bytes []byte) (IKeg, error) {
	k := &pbKeg.Keg{}
	if err := proto.Unmarshal(bytes, k); err != nil {
		return nil, err
	}

	return FromProto(k), nil
}

// FromProto converts a proto keg to a model.Keg
func FromProto(k *pbKeg.Keg) IKeg {
	var tree merkle.ITree
	if k.GetTree() != nil {
		tree = merkle.FromProto(k.GetTree(), wrapper)
	} else {
		tree = merkle.NewTree(16)
	}

	return &Keg{
		id:                 k.Id,
		options:            optionsFromProto(k.Options),
		liquidByAccessName: make(map[string]string),
		liquidInfo:         make(map[string]liquid.IInfo),
		merkleTree:         tree,
		lastUpdated:        k.LastUpdated,
		deleted:            false,
	}
}

func wrapper() merkle.IContent {
	return liquid.NewEmptyMerkleTreeLiquid()
}

// FromDir tries to load a keg form a FS dir
func FromDir(directory string) (IKeg, error) {
	log.Printf("loading keg %s", directory)
	var content []byte
	var err error
	var k = NewEmptyKeg()

	localKegFile := fmt.Sprintf("%s/%s/%s", config.C.DataRoot, directory, kegFile)
	if content, err = ioutil.ReadFile(localKegFile); err != nil {
		log.Printf("error loading keg file /%s/%s", directory, kegFile)
		return nil, err
	}

	if k, err = FromBytes(content); err != nil {
		log.Printf("corrupted keg file /%s/%s", directory, kegFile)
		return nil, err
	}

	files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", config.C.DataRoot, directory))
	if err != nil {
		log.Printf("error reading directory /%s", directory)
		return k, nil
	}

	for _, file := range files {
		if !file.IsDir() {
			if l, err := liquid.FromFile(fmt.Sprintf("%s/%s/%s", config.C.DataRoot, directory, file.Name())); err == nil {
				info := l.GetLiquidInfo()
				k.AddLiquid(info)
			} else {
				log.Println(err)
			}
		}
	}

	log.Printf("loaded keg %s\n", k.GetID())
	return k, nil
}
