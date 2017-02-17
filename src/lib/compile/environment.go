package compile

import "../vm"

type environment struct {
	parent Gettable
	me     map[string]*vm.Thunk
}

// TODO: We may not need to nest environments because closures are removed and
// normalized in IR.
func newEnvironment(parent *environment) *environment {
	g := Gettable(parent)

	if parent == nil {
		g = prelude
	}

	return &environment{parent: g, me: make(map[string]*vm.Thunk)}
}

func (e *environment) set(s string, t *vm.Thunk) {
	e.me[s] = t
}

func (e *environment) get(s string) *vm.Thunk {
	t, ok := e.me[s]

	if ok {
		return t
	}

	return e.parent.get(s)
}
