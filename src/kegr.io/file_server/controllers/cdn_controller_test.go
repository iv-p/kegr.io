package controllers

import (
	"cerberus/config"
	"cerberus/service"
	"testing"
)

type fakeKegService struct {
	service.IKegService
}

func TestCreateNewCdnController(t *testing.T) {
	c := &config.Config{}
	ks := &fakeKegService{}

	cc := NewCdnController(c, ks)

	if cc.kegService != ks {
		t.Error("keg service mismatch")
	}
	if cc.c != c {
		t.Error("config mismatch")
	}
}
