package compile

import (
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/raviqqe/tisp/src/lib/util"
)

type environment struct {
	me       map[string]*core.Thunk
	fallback func(string) (*core.Thunk, error)
}

func newEnvironment(fallback func(string) (*core.Thunk, error)) environment {
	return environment{
		me:       make(map[string]*core.Thunk),
		fallback: fallback,
	}
}

func (e *environment) set(s string, t *core.Thunk) {
	e.me[s] = t
}

func (e environment) get(s string) *core.Thunk {
	if t, ok := e.me[s]; ok {
		return t
	}

	t, err := e.fallback(s)

	if err == nil {
		return t
	}

	util.Fail("The name, %s is not found.", s)
	panic("unreachable")
}
