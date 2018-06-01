package merkle

import (
	"bytes"
	"testing"
)

const (
	leafTestItems = 16
)

func TestLeafAddHashEqual(t *testing.T) {
	one := NewLeaf()
	two := NewLeaf()

	for i := 0; i < 3; i++ {
		c := newC()
		one.Add(c)
		two.Add(c)
	}

	if !bytes.Equal(one.GetHash(), two.GetHash()) {
		t.Error("hashes should be equal")
	}
}

func TestLeafAddHashDifferent(t *testing.T) {
	one := NewLeaf()
	two := NewLeaf()

	for i := 0; i < leafTestItems; i++ {
		c := newC()
		v := newC()
		one.Add(c)
		two.Add(v)
	}

	if bytes.Equal(one.GetHash(), two.GetHash()) {
		t.Error("hashes should be equal")
	}
}

func TestLeafDeleteHash(t *testing.T) {
	one := NewLeaf()

	hash := one.GetHash()
	for i := 0; i < leafTestItems; i++ {
		c := newC()
		one.Add(c)
		one.Delete(c.GetID())
	}

	if !bytes.Equal(hash, one.GetHash()) {
		t.Error("hashes should be equal")
	}
}
func TestLeafHash(t *testing.T) {
	one := NewLeaf()
	two := NewLeaf()

	for i := 0; i < leafTestItems; i++ {
		c := newC()
		v := newC()

		one.Add(c)
		one.Delete(c.GetID())

		two.Add(v)
		two.Delete(v.GetID())
	}

	if !bytes.Equal(one.GetHash(), two.GetHash()) {
		t.Error("hashes should be equal")
	}
}
