package merkle

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"testing"
)

const (
	treeTestItems = 1024
	treeDepth     = 16
)

type C struct {
	id          []byte
	lastUpdated int64
	Content
}

func newC() *C {
	id := make([]byte, 8)
	rand.Read(id)
	return &C{
		id: id,
	}
}

func (c *C) GetID() []byte {
	return c.id
}

func (c *C) SetLastUpdated(lastUpdated int64) {
	c.lastUpdated = lastUpdated
}

func (c *C) GetHash() []byte {
	hash := sha1.New()
	hash.Write(c.id)
	return hash.Sum(nil)
}

func setup() *Tree {
	return NewTree(treeDepth)
}

func TestTreeAdd(t *testing.T) {
	tree := setup()
	hash := tree.Hash()

	for i := 0; i < treeTestItems; i++ {
		c := newC()
		if err := tree.Add(c); err != nil {
			t.Error(err)
		}

		if err := tree.Delete(c.GetID()); err != nil {
			t.Error(err)
		}
	}

	newHash := tree.Hash()
	if bytes.Compare(hash, newHash) != 0 {
		t.Error("Hash not the same")
	}
}

func TestTreeDiffEqual(t *testing.T) {
	one := setup()
	two := setup()

	for i := 0; i < treeTestItems; i++ {
		c := newC()
		if err := one.Add(c); err != nil {
			t.Error(err)
		}
		if err := two.Add(c); err != nil {
			t.Error(err)
		}
	}

	if _, eq := one.Diff(two); !eq {
		t.Error("trees should be equal")
	}
}

func TestTreeDiffNotEqualMakeEqual(t *testing.T) {
	one := setup()
	two := setup()

	for i := 0; i < treeTestItems; i++ {
		c := newC()
		if err := two.Add(c); err != nil {
			t.Error(err)
		}
	}

	var d []Content
	var eq bool

	if d, eq = one.Diff(two); eq {
		t.Error("trees should not be equal")
	}

	for _, item := range d {
		one.Add(item)
	}

	if d, eq = one.Diff(two); !eq {
		t.Error("trees should be equal")
	}
}

func TestTreeDiffNotEqualAllElements(t *testing.T) {
	one := setup()
	two := setup()

	for i := 0; i < treeTestItems; i++ {
		if err := one.Add(newC()); err != nil {
			t.Error(err)
		}
		if err := two.Add(newC()); err != nil {
			t.Error(err)
		}
	}

	var eq bool

	if _, eq = one.Diff(two); eq {
		t.Error("trees should not be equal")
	}
}

func TestTreeDiffNotEqualOnlyOneElement(t *testing.T) {
	one := setup()
	two := setup()

	for i := 0; i < treeTestItems; i++ {
		c := newC()
		if err := one.Add(c); err != nil {
			t.Error(err)
		}
		if err := two.Add(c); err != nil {
			t.Error(err)
		}
	}

	two.Add(newC())

	var eq bool

	if _, eq = one.Diff(two); eq {
		t.Error("trees should not be equal")
	}
}

func TestTreeSameHashWhenCreated(t *testing.T) {
	one := setup()
	two := setup()

	if !bytes.Equal(one.Hash(), two.Hash()) {
		t.Error("hashes not equal")
	}
}

func TestTreeSameHashAfterUpdate(t *testing.T) {
	one := setup()
	two := setup()

	c := newC()

	one.Add(c)
	c.SetLastUpdated(16)
	one.Add(c)

	two.Add(c)

	if !bytes.Equal(one.Hash(), two.Hash()) {
		t.Error("hashes not equal")
	}
}
