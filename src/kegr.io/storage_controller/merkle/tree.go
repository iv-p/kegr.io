package merkle

import (
	"bytes"
	"errors"

	pbMerkle "kegr.io/protobuf/model/merkle"
	"kegr.io/storage_controller/util"
)

// Tree is the main structure for merkle trees
type Tree struct {
	rootNode *node
	depth    int
}

// ITree is the tree interface
type ITree interface {
	Add(IContent) error
	Update(IContent) error
	Delete([]byte) error
	Hash() []byte
	Diff(ITree) []IContent
	ToProto() *pbMerkle.Tree

	getNode() *node
}

// NewTree initialises a new Merkle tree and returns the object
func NewTree(depth int) ITree {
	t := &Tree{
		depth:    depth,
		rootNode: &node{},
	}

	return t
}

// Add adds a content object to the Merkle tree
func (t *Tree) Add(content IContent) error {
	id := content.GetID()
	var err error

	current := t.rootNode
	for depth := 0; depth < t.depth; depth++ {
		var b byte
		if b, err = util.GetBitFromByteArray(depth, id); err != nil {
			return errors.New("Error inserting")
		}

		if b == 1 {
			if current.left == nil {
				current.left = &node{}
				current.left.parent = current
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &node{}
				current.right.parent = current
			}
			current = current.right
		}
	}

	if current.leaf == nil {
		current.leaf = NewLeaf()
	}

	current.leaf.Add(content)
	return nil
}

// Update a node in the tree
func (t *Tree) Update(content IContent) error {
	var err error
	notFound := errors.New("ID not found in merkle tree")

	current := t.rootNode
	for depth := 0; depth < t.depth; depth++ {
		var b byte
		if b, err = util.GetBitFromByteArray(depth, content.GetID()); err != nil {
			return errors.New("Error deleting")
		}

		if b == 1 {
			if current.left == nil {
				return notFound
			}
			current = current.left
		} else {
			if current.right == nil {
				return notFound
			}
			current = current.right
		}
	}

	if current.leaf == nil {
		return notFound
	}

	if _, exist := current.leaf.Get(content.GetID()); !exist {
		return notFound
	}

	current.leaf.Add(content)
	return nil
}

// Delete removes a content object to the Merkle tree
func (t *Tree) Delete(id []byte) error {
	var err error
	notFound := errors.New("ID not found in merkle tree")

	current := t.rootNode
	for depth := 0; depth < t.depth; depth++ {
		var b byte
		if b, err = util.GetBitFromByteArray(depth, id); err != nil {
			return errors.New("Error deleting")
		}

		if b == 1 {
			if current.left == nil {
				return notFound
			}
			current = current.left
		} else {
			if current.right == nil {
				return notFound
			}
			current = current.right
		}
	}

	if current.leaf == nil {
		return notFound
	}

	current.leaf.Delete(id)
	return t.cleanupAfterDelete(id, current)
}

// Diff takes another merkle tree and traverses both to find differences.
// It returns the additions and deletion differences
func (t *Tree) Diff(other ITree) []IContent {
	if bytes.Equal(t.Hash(), other.Hash()) {
		return nil
	}

	return t.rootNode.diff(other.getNode())
}

// Hash returns the root hash of the merkle tree
func (t *Tree) Hash() []byte {
	return t.rootNode.getHash()
}

// ToProto returns the protobuf representation of the tree
func (t *Tree) ToProto() *pbMerkle.Tree {
	return &pbMerkle.Tree{
		Depth:    int64(t.depth),
		RootNode: t.rootNode.toProto(),
	}
}

// FromProto reads a proto tree and returns a tree structure
func FromProto(t *pbMerkle.Tree, newContentObject func() IContent) *Tree {
	return &Tree{
		depth:    int(t.Depth),
		rootNode: nodeFromProto(t.GetRootNode(), newContentObject),
	}
}

func (t *Tree) cleanupAfterDelete(id []byte, node *node) error {

	if !node.leaf.isEmpty() {
		return nil
	}

	depth := t.depth - 1
	current := node.parent
	node.leaf = nil
	var err error
	var b byte

	for current != nil {
		if b, err = util.GetBitFromByteArray(depth, id); err != nil {
			return errors.New("Error deleting")
		}
		deleted := false

		if b == 1 {
			if current.left.left == nil && current.left.right == nil {
				deleted = true
				current.left = nil
			}
		} else {
			if current.right.left == nil && current.right.right == nil {
				deleted = true
				current.right = nil
			}
		}

		if !deleted {
			break
		}
		current = current.parent
		depth--
	}

	return nil
}

func (t *Tree) getNode() *node {
	return t.rootNode
}
