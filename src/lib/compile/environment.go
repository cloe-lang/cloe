package compile

import (
	"fmt"

	"github.com/cloe-lang/cloe/src/lib/core"
)

type environment struct {
	me       module
	fallback func(string) (core.Value, error)
}

func newEnvironment(fallback func(string) (core.Value, error)) environment {
	return environment{module{}, fallback}
}

func (e *environment) set(s string, t core.Value) {
	e.me[s] = t
}

func (e environment) get(s string) (core.Value, error) {
	if v, ok := e.me[s]; ok {
		return v, nil
	}

	v, err := e.fallback(s)

	if err == nil {
		return v, nil
	}

	return nil, fmt.Errorf("name, %s not found", s)
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
