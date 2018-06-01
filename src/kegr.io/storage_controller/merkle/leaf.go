package merkle

import (
	"bytes"
	"crypto/md5"
	"sort"

	pbMerkle "kegr.io/protobuf/model/merkle"
)

// Leaf is the struct that holds the actual items in the merkle tree
type Leaf struct {
	content map[string]IContent
	hash    []byte
}

type IContent interface {
	GetID() []byte
	SetID(id []byte)
	GetHash() []byte
	SetHash(hash []byte)
	GetLastUpdated() int64
	SetLastUpdated(lastUpdated int64)
	GetProto() *pbMerkle.Content
}

// NewLeaf returns an initialised Lead object
func NewLeaf() *Leaf {
	return &Leaf{
		content: make(map[string]IContent),
	}
}

// GetHash returns the has of the contents of the leaf
func (l *Leaf) GetHash() []byte {
	var sortedIds []string
	for id := range l.content {
		sortedIds = append(sortedIds, id)
	}
	sort.Strings(sortedIds)

	hash := md5.New()
	for _, key := range sortedIds {
		hash.Write(l.content[key].GetHash())
	}
	return hash.Sum(nil)
}

// Add inserts a new item in the leaf's content
func (l *Leaf) Add(item IContent) {
	l.content[string(item.GetID())] = item
}

// GetContent returns the content of the leaf
func (l *Leaf) GetContent() []IContent {
	var content []IContent
	for _, v := range l.content {
		content = append(content, v)
	}
	return content
}

// Get returns the content item with that id and a bool reporesenting
// if that item is present
func (l *Leaf) Get(id []byte) (IContent, bool) {
	if value, exist := l.content[string(id)]; exist {
		return value, true
	}
	return nil, false
}

// Delete removes an item from the content if found
func (l *Leaf) Delete(id []byte) {
	if _, exist := l.content[string(id)]; !exist {
		return
	}
	delete(l.content, string(id))
}

func (l *Leaf) isEmpty() bool {
	return len(l.content) == 0
}

func (l *Leaf) diff(other *Leaf) []IContent {
	if bytes.Equal(l.GetHash(), other.GetHash()) {
		return nil
	}
	var diff []IContent

	for lk, lv := range l.content {
		if ov, exist := other.content[lk]; exist {
			if !bytes.Equal(lv.GetHash(), ov.GetHash()) && ov.GetLastUpdated() > lv.GetLastUpdated() {
				diff = append(diff, ov)
			}
		}
	}

	for ok, ov := range other.content {
		if _, exist := l.content[ok]; !exist {
			diff = append(diff, ov)
		}
	}

	return diff
}

func (l *Leaf) toProto() *pbMerkle.Leaf {
	pbLeaf := &pbMerkle.Leaf{
		Hash: l.hash,
	}

	pbContent := make(map[string]*pbMerkle.Content)

	for k, v := range l.content {
		pbContent[k] = v.GetProto()
	}

	pbLeaf.Content = pbContent

	return pbLeaf
}

func leafFromProto(l *pbMerkle.Leaf, newContentObject func() IContent) *Leaf {
	ll := &Leaf{
		hash: l.GetHash(),
	}

	content := make(map[string]IContent)

	for k, v := range l.GetContent() {
		content[k] = newContentObject()
		content[k].SetID(v.GetID())
		content[k].SetHash(v.GetHash())
		content[k].SetLastUpdated(v.GetLastUpdated())
	}

	ll.content = content

	return ll
}
