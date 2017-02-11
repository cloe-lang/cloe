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

func (s Set) Remove(x interface{}) Set {
	return Set{s.Tree.Remove(x)}
}

func (s Set) FirstRest() (interface{}, Set) {
	x, t := s.Tree.FirstRest()
	return x, Set{t}
}

func (s Set) Merge(ss Set) Set {
	return Set{s.Tree.Merge(ss.Tree)}
}
