package store

import (
	"slices"
)

type node struct {
	keys     []int
	children []*node
}

type BTree struct {
	root *node
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

	return insertToNode(node, parent, key)
}

func insertToNode(node *node, parent *node, key int) error {
	operation := getInsertOperation(node)
	switch operation {
	case InsertDirecty:
		insertDirectly(node, key)
	case Split:
		split(node, parent, key)
	}
	return nil
}

func findNodeForInsert(node *node, parent *node, key int) (*node, *node, error) {

	if node != nil && len(node.children) > 0 {
		child := findChild(node, key)
		return findNodeForInsert(child, node, key)
	}

	return node, parent, nil
}

func findChild(node *node, key int) *node {
	for i, childKey := range node.keys {
		if key < childKey {
			return node.children[i]
		}
	}

	return node.children[len(node.children)-1]
}

func newNode() *node {
	return &node{keys: []int{}}
}

func getInsertOperation(node *node) Operation {
	return InsertDirecty
}

func insertDirectly(node *node, key int) {
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

func split(node *node, parent *node, key int) {
	allKeys := append(node.keys, key)
	slices.Sort(allKeys)
	mid := len(allKeys) / 2
	leftKeys := allKeys[:mid]
	middleKey := allKeys[mid]
	rightKeys := allKeys[mid+1:]

	// This node gets the left half of the keys:
	node.keys = leftKeys

	// The middle key is pushed up to the parent:
	insertToNode(node, parent, middleKey)

	// The right keys are moved to a new right sibling:
	createRightSibling(node, parent, rightKeys)
}

func createRightSibling(node *node, parent *node, keys []int) {
	sibling := newNode()
	sibling.keys = keys
	nodeIndex := slices.Index(parent.children, node)
	parent.children = slices.Insert(parent.children, nodeIndex+1, sibling)
}
