package ir

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/core"
)

// Switch represents a switch expression.
type Switch struct {
	value interface{}
	cases []Case
}

// NewSwitch creates a switch expression.
func NewSwitch(v interface{}, cs []Case) Switch {
	if len(cs) == 0 {
		panic(fmt.Errorf("A number of cases in switch expressions must be more than 0"))
	}

	return Switch{v, cs}
}

func (s Switch) compileToDict() *core.Thunk {
	ks := make([]core.Value, 0, len(s.cases))
	vs := make([]*core.Thunk, 0, len(s.cases))

	for _, c := range s.cases {
		ks = append(ks, c.pattern.Eval())
		vs = append(vs, core.NewNumber(float64(c.value)))
	}

	return core.NewDictionary(ks, vs)
}
