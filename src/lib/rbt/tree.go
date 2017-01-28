package rbt

type Tree struct {
	node *node
	less func(interface{}, interface{}) bool
}

func NewTree(less func(interface{}, interface{}) bool) Tree {
	return Tree{
		node: nil,
		less: less,
	}
}

func (t Tree) Insert(x interface{}) Tree {
	return Tree{
		node: t.node.insert(x, t.less),
		less: t.less,
	}
}

func (t Tree) Search(x interface{}) (interface{}, bool) {
	return t.node.search(x, t.less)
}

func (t Tree) Remove(x interface{}) (Tree, bool) {
	n, ok := t.node.remove(x, t.less)

	return Tree{
		node: n,
		less: t.less,
	}, ok
}
