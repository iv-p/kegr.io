package keg

import "kegr.io/storage_controller/merkle"

type KegDiff struct {
	Options IOptions
	Content []merkle.IContent
}

func NewKegDiff() *KegDiff {
	return &KegDiff{}
}
