package ir

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/core"
)

// Switch represents a switch expression.
type Switch struct {
	value       interface{}
	cases       []Case
	defaultCase interface{}
}

// NewSwitch creates a switch expression.
func NewSwitch(v interface{}, cs []Case, d interface{}) Switch {
	if len(cs) == 0 && d == nil {
		panic(fmt.Errorf("A number of cases in switch expressions must be more than 0"))
	}

	return Switch{v, cs, d}
}

func (s Switch) compileToDict() *core.Thunk {
	ks := make([]core.Value, 0, len(s.cases))
	vs := make([]*core.Thunk, 0, len(s.cases))

	for i, c := range s.cases {
		ks = append(ks, c.pattern.Eval())
		vs = append(vs, core.NewNumber(float64(i)))
	}

	return core.NewDictionary(ks, vs)
}
