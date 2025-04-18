package store

import (
	"slices"
)

type Node struct {
	keys     []int
	children []*Node
}

type BTree struct {
	root *Node
}

type Operation int

const (
	InsertDirecty Operation = iota
	Split
)

func (tree *BTree) Insert(key int) error {
	if tree.root == nil {
		tree.root = newNode()
	}

	node, parent, err := findNodeForInsert(tree.root, nil, key)
	if err != nil {
		return err
	}

	operation := getInsertOperation(node)
	switch operation {
	case InsertDirecty:
		insertDirectly(node, key)
	case Split:
		if parent == nil {
			// This is a root split
			splitRoot(tree, node, key)
		} else {
			split(node, parent, key)
		}
	}
	return nil
}

func findNodeForInsert(node *Node, parent *Node, key int) (*Node, *Node, error) {
	if node != nil && len(node.children) > 0 {
		child := findChild(node, key)
		return findNodeForInsert(child, node, key)
	}

	return node, parent, nil
}

func findChild(node *Node, key int) *Node {
	for i, childKey := range node.keys {
		if key < childKey {
			return node.children[i]
		}
	}

	return node.children[len(node.children)-1]
}

func newNode() *Node {
	return &Node{keys: []int{}}
}

func getInsertOperation(node *Node) Operation {
	if len(node.keys) >= 4 { // Maximum 4 keys per node
		return Split
	}
	return InsertDirecty
}

func insertDirectly(node *Node, key int) {
	length := len(node.keys)
	index := length
	for i := 0; i < length; i++ {
		if key < node.keys[i] {
			index = i
			break
		}
	}
	node.keys = slices.Insert(node.keys, index, key)
}

func splitRoot(tree *BTree, node *Node, key int) {
	allKeys := append(node.keys, key)
	slices.Sort(allKeys)
	mid := len(allKeys) / 2
	leftKeys := allKeys[:mid]
	middleKey := allKeys[mid]
	rightKeys := allKeys[mid+1:]

	// Create left and right children
	leftChild := newNode()
	leftChild.keys = leftKeys
	rightChild := newNode()
	rightChild.keys = rightKeys

	// Create new root
	newRoot := newNode()
	newRoot.keys = []int{middleKey}
	newRoot.children = []*Node{leftChild, rightChild}

	// Update the tree's root
	tree.root = newRoot
}

func split(node *Node, parent *Node, key int) {
	allKeys := append(node.keys, key)
	slices.Sort(allKeys)
	mid := len(allKeys) / 2
	leftKeys := allKeys[:mid]
	middleKey := allKeys[mid]
	rightKeys := allKeys[mid+1:]

	// This node gets the left half of the keys
	node.keys = leftKeys

	// Create new right sibling
	rightSibling := newNode()
	rightSibling.keys = rightKeys

	// Insert the middle key into the parent
	insertDirectly(parent, middleKey)

	// Find the index of the current node in parent's children
	nodeIndex := slices.Index(parent.children, node)

	// Insert the right sibling after the current node
	parent.children = slices.Insert(parent.children, nodeIndex+1, rightSibling)
}
