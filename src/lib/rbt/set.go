package rbt

type Set struct{ Tree }

func NewSet(less func(interface{}, interface{}) bool) Set {
	return Set{NewTree(less)}
}

func (s Set) Insert(x interface{}) Set {
	return Set{s.Tree.Insert(x)}
}

func (s Set) Include(x interface{}) bool {
	_, ok := s.Tree.Search(x)
	return ok
}

func (s Set) Remove(x interface{}) (Set, bool) {
	t, ok := s.Tree.Remove(x)
	return Set{t}, ok
}
