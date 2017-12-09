package compile

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/core"
)

type environment struct {
	me       module
	fallback func(string) (*core.Thunk, error)
}

func newEnvironment(fallback func(string) (*core.Thunk, error)) environment {
	return environment{
		me:       make(module),
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

	panic(fmt.Errorf("The name, %s is not found", s))
}

func (e environment) toMap() module {
	return e.me
}

func (e environment) copy() environment {
	m := make(module, len(e.me))

	for k, v := range e.me {
		m[k] = v
	}

	return environment{m, e.fallback}
}
