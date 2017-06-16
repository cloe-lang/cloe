package rbt

// Set represents a set of values.
type Set struct{ Tree }

// NewSet creates an empty Set from a function which determines orders of values inside.
func NewSet(compare func(interface{}, interface{}) int) Set {
	return Set{NewTree(compare)}
}

// Insert inserts a value to a set.
func (s Set) Insert(x interface{}) Set {
	return Set{s.Tree.Insert(x)}
}

// Include returns true if a value is included in a set, or false otherwise.
func (s Set) Include(x interface{}) bool {
	_, ok := s.Tree.Search(x)
	return ok
}

// Remove creates a set without a value.
func (s Set) Remove(x interface{}) Set {
	return Set{s.Tree.Remove(x)}
}

// FirstRest returns a value which was in the original set and a new set without it.
func (s Set) FirstRest() (interface{}, Set) {
	x, t := s.Tree.FirstRest()
	return x, Set{t}
}

// Merge merges 2 sets.
func (s Set) Merge(ss Set) Set {
	return Set{s.Tree.Merge(ss.Tree)}
}
