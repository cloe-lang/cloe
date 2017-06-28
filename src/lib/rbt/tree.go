package rbt

// Tree represents an red-black tree.
type Tree struct {
	node    *node
	compare func(interface{}, interface{}) int
}

// NewTree creates a empty red-black tree.
// The compare function should return a negative when x < y, 0 when x = y, and a positive when x > y.
func NewTree(compare func(interface{}, interface{}) int) Tree {
	return Tree{
		node:    nil,
		compare: compare,
	}
}

// Insert inserts an element into a tree.
func (t Tree) Insert(x interface{}) Tree {
	return Tree{
		node:    t.node.insert(x, t.compare),
		compare: t.compare,
	}
}

// Search searches an element in a tree.
// It returns a found element in addition to a condition if the element is
// found because a compare function passed to NewTree can compare elements
// partially.
func (t Tree) Search(x interface{}) (interface{}, bool) {
	return t.node.search(x, t.compare)
}

// Remove removes an element in a tree.
func (t Tree) Remove(x interface{}) Tree {
	return Tree{
		node:    t.node.remove(x, t.compare),
		compare: t.compare,
	}
}

// FirstRest returns an element inside and a tree without it.
func (t Tree) FirstRest() (interface{}, Tree) {
	x := t.node.min()

	if x == nil {
		return nil, NewTree(t.compare)
	}

	return x, t.Remove(x)
}

// Empty returns true when the tree has no element or false otherwise.
func (t Tree) Empty() bool {
	return t.node == nil
}

// Size returns a number of elements inside a tree.
func (t Tree) Size() int {
	return t.node.size()
}

// Merge merges 2 trees and returns a merged trees.
func (t Tree) Merge(merged Tree) Tree {
	for {
		var x interface{}
		x, merged = merged.FirstRest()

		if x == nil {
			break
		}

		t = t.Insert(x)
	}

	return t
}

// Dump prints a tree to stdout.
func (t Tree) Dump() {
	t.node.dump()
}
