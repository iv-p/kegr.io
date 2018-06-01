package merkle

import (
	"bytes"
	"crypto/md5"
	"log"

	pbMerkle "kegr.io/protobuf/model/merkle"
)

// Node is the stucture that defines a single node in the
// merkle tree
type node struct {
	parent *node
	left   *node
	right  *node
	hash   []byte
	leaf   *Leaf
}

// Hash computes the sum of the hashes of the node's children
func (n *node) getHash() []byte {
	hash := md5.New()
	if n.left == nil && n.right == nil && n.leaf != nil {
		hash.Write(n.leaf.GetHash())
	} else {
		if n.left != nil {
			hash.Write(n.left.getHash())
		}
		if n.right != nil {
			hash.Write(n.right.getHash())
		}
	}
	return hash.Sum(nil)
}

func (n *node) recomputeHash() {
	hash := md5.New()
	if n.left == nil && n.right == nil && n.leaf != nil {
		hash.Write(n.leaf.GetHash())
	} else {
		if n.left != nil {
			hash.Write(n.left.getHash())
		}
		if n.right != nil {
			hash.Write(n.right.getHash())
		}
	}
	n.hash = hash.Sum(nil)
}

// ID dummy function to fulfil interface. NOT USED
func (n *node) id() []byte { return nil }

func (n *node) isLeaf() bool {
	return n.leaf != nil && n.left == nil && n.right == nil
}

// diff calculates the diff between two nodes, returning additions and deletions
// if two separate lists
func (n *node) diff(other *node) []IContent {
	var diff []IContent

	if other == nil {
		return nil
	}

	if bytes.Equal(n.getHash(), other.getHash()) {
		return nil
	}

	if n.isLeaf() && other.isLeaf() {
		return n.leaf.diff(other.leaf)
	} else if n.isLeaf() || other.isLeaf() {
		log.Println("merkle tree inconsistency")
	}

	if n.left == nil && other.left != nil {
		diff = append(diff, other.left.getAllContent()...)
	}

	if n.left != nil && other.left != nil {
		diff = append(diff, n.left.diff(other.left)...)
	}

	if n.right == nil && other.right != nil {
		diff = append(diff, other.right.getAllContent()...)
	}

	if n.right != nil && other.right != nil {
		diff = append(diff, n.right.diff(other.right)...)
	}

	return diff
}

func (n *node) getAllContent() []IContent {
	var content []IContent
	if n.leaf != nil {
		return n.leaf.GetContent()
	}

	if n.left != nil {
		content = append(content, n.left.getAllContent()...)
	}

	if n.right != nil {
		content = append(content, n.right.getAllContent()...)
	}

	return content
}

func (n *node) toProto() *pbMerkle.Node {
	pbNode := &pbMerkle.Node{
		Hash: n.getHash(),
	}
	if n.left != nil {
		pbNode.Left = n.left.toProto()
	}

	if n.right != nil {
		pbNode.Right = n.right.toProto()
	}

	if n.leaf != nil {
		pbNode.Leaf = n.leaf.toProto()
	}

	return pbNode
}

func nodeFromProto(n *pbMerkle.Node, newContentObject func() IContent) *node {
	nn := &node{
		hash: n.GetHash(),
	}

	if n.GetLeft() != nil {
		nn.left = nodeFromProto(n.GetLeft(), newContentObject)
	}

	if n.GetRight() != nil {
		nn.right = nodeFromProto(n.GetRight(), newContentObject)
	}

	if n.GetLeaf() != nil {
		nn.leaf = leafFromProto(n.GetLeaf(), newContentObject)
	}

	return nn
}
