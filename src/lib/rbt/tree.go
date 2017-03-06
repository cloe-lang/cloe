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

func (t Tree) Remove(x interface{}) Tree {
	return Tree{
		node: t.node.remove(x, t.less),
		less: t.less,
	}
}

func (t Tree) FirstRest() (interface{}, Tree) {
	x := t.node.min()

	if x == nil {
		return nil, NewTree(t.less)
	}

	return x, t.Remove(x)
}

func (t Tree) Empty() bool {
	return t.node == nil
}

func (t Tree) Size() int {
	return t.node.size()
}

func (t Tree) Merge(tt Tree) Tree {
	for {
		var x interface{}
		x, tt = tt.FirstRest()

		if x == nil {
			break
		}

		t = t.Insert(x)
	}

	return t
}

func (t Tree) Dump() {
	t.node.dump()
}
