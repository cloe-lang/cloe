package env

import (
	"../../vm"
	"fmt"
)

type Environment struct {
	parent   *Environment
	me       map[string]*vm.Thunk
	fallback func(string) (*vm.Thunk, error)
}

func NewEnvironment(fallback func(string) (*vm.Thunk, error)) Environment {
	return Environment{
		parent:   nil,
		me:       make(map[string]*vm.Thunk),
		fallback: fallback,
	}
}

func (e *Environment) Set(s string, t *vm.Thunk) {
	e.me[s] = t
}

func (e Environment) Get(s string) (*vm.Thunk, error) {
	if t, ok := e.me[s]; ok {
		return t, nil
	}

	if e.parent != nil {
		return e.parent.Get(s)
	}

	return e.fallback(s)
}

func (e Environment) Child() Environment {
	return Environment{
		parent:   &e,
		me:       make(map[string]*vm.Thunk),
		fallback: e.fallback,
	}
}

func (e Environment) Parent() Environment {
	if e.parent == nil {
		panic(fmt.Errorf("Parent is nil."))
	}

	return *e.parent
}
