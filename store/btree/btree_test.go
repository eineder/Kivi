package store

import (
	"fmt"
	"testing"
)

func TestBTreeInsert(t *testing.T) {
	tests := []struct {
		name  string
		given *Node
		when  int
		then  *Node
	}{
		{
			name:  "Empty tree",
			given: nil,
			when:  3,
			then:  &Node{keys: []int{3}},
		},
		{
			name:  "Insert to root",
			given: &Node{keys: []int{1, 2, 4}},
			when:  3,
			then:  &Node{keys: []int{1, 2, 3, 4}},
		},
		{
			name:  "Split root",
			given: &Node{keys: []int{1, 2, 4, 5}},
			when:  3,
			then: &Node{
				keys: []int{3},
				children: []*Node{
					{keys: []int{1, 2}},
					{keys: []int{4, 5}},
				},
			},
		},
		{
			name: "Insert to leaf (middle)",
			given: &Node{
				keys:     []int{10, 20, 30, 40},
				children: []*Node{{keys: []int{5, 6, 7}}, {keys: []int{15, 16, 19}}},
			},
			when: 17,
			then: &Node{
				keys:     []int{10, 20, 30, 40},
				children: []*Node{{keys: []int{5, 6, 7}}, {keys: []int{15, 16, 17, 19}}},
			},
		},
		{
			name: "Insert to leaf (start)",
			given: &Node{
				keys:     []int{10, 20, 30, 40},
				children: []*Node{{keys: []int{5, 6, 7}}, {keys: []int{15, 16, 19}}},
			},
			when: 14,
			then: &Node{
				keys:     []int{10, 20, 30, 40},
				children: []*Node{{keys: []int{5, 6, 7}}, {keys: []int{14, 15, 16, 19}}},
			},
		},
		{
			name: "Insert to leaf (end)",
			given: &Node{
				keys:     []int{10, 20, 30, 40},
				children: []*Node{{keys: []int{5, 6, 7}}, {keys: []int{15, 16, 18}}},
			},
			when: 19,
			then: &Node{
				keys:     []int{10, 20, 30, 40},
				children: []*Node{{keys: []int{5, 6, 7}}, {keys: []int{15, 16, 18, 19}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := &BTree{root: tt.given}
			tree.Insert(tt.when)
			// Compare the tree with the expected tree
			equals, msg := Equals(tree.root, tt.then)
			if !equals {
				t.Error(msg)
			}
		})
	}
}

// func TestFindSliceInsert(t *testing.T) {
// 	sl := []int{3, 6, 7}
// 	new := slices.Insert(sl, 3, 333)

// 	t.Errorf("%+v", new)
// }

func TestFindChild(t *testing.T) {
	tests := []struct {
		name     string
		node     *Node
		key      int
		expected *Node
	}{
		{
			name: "Find child in the middle",
			node: &Node{
				keys:     []int{10, 20, 30},
				children: []*Node{{keys: []int{5}}, {keys: []int{15}}, {keys: []int{25}}, {keys: []int{35}}},
			},
			key:      22,
			expected: &Node{keys: []int{25}},
		},
		{
			name: "Find child at the beginning",
			node: &Node{
				keys:     []int{10, 20, 30},
				children: []*Node{{keys: []int{5}}, {keys: []int{15}}, {keys: []int{25}}, {keys: []int{35}}},
			},
			key:      8,
			expected: &Node{keys: []int{5}},
		},
		{
			name: "Find child at the end",
			node: &Node{
				keys:     []int{10, 20, 30},
				children: []*Node{{keys: []int{5}}, {keys: []int{15}}, {keys: []int{25}}, {keys: []int{35}}},
			},
			key:      32,
			expected: &Node{keys: []int{35}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			child := findChild(tt.node, tt.key)
			equals, msg := Equals(child, tt.expected)
			if !equals {
				t.Error(msg)
			}
		})
	}
}

func Equals(node1 *Node, node2 *Node) (bool, string) {
	if node1 == nil && node2 == nil {
		return true, ""
	}
	if node1 == nil || node2 == nil {
		return false, "one node is nil while the other is not"
	}

	node1keysCount := len(node1.keys)
	node2keysCount := len(node2.keys)

	if node1keysCount != node2keysCount {
		return false, fmt.Sprintf("node1 has %d keys, node2 has %d keys", node1keysCount, node2keysCount)
	}

	for i := 0; i < len(node1.keys); i++ {
		node1key := node1.keys[i]
		node2key := node2.keys[i]
		if node1key != node2key {
			return false, fmt.Sprintf("node1 has key %d and node2 has key %d at index %d", node1key, node2key, i)
		}
	}

	if len(node1.children) != len(node2.children) {
		return false, fmt.Sprintf("node1 has %d children, node2 has %d children", len(node1.children), len(node2.children))
	}

	for i, child1 := range node1.children {
		child2 := node2.children[i]
		equal, msg := Equals(child1, child2)
		if !equal {
			return false, msg
		}
	}

	return true, ""
}
