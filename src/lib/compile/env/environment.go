package env

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/core"
)

type Environment struct {
	parent   *Environment
	me       map[string]*core.Thunk
	fallback func(string) (*core.Thunk, error)
}

func NewEnvironment(fallback func(string) (*core.Thunk, error)) Environment {
	return Environment{
		parent:   nil,
		me:       make(map[string]*core.Thunk),
		fallback: fallback,
	}
}

func (e *Environment) Set(s string, t *core.Thunk) {
	e.me[s] = t
}

func (e Environment) Get(s string) (*core.Thunk, error) {
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
		me:       make(map[string]*core.Thunk),
		fallback: e.fallback,
	}
}

func (e Environment) Parent() Environment {
	if e.parent == nil {
		panic(fmt.Errorf("parent environment is nil"))
	}

	return *e.parent
}
