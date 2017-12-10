package ir

import (
	"github.com/tisp-lang/tisp/src/lib/core"
)

// Switch represents a switch expression.
type Switch struct {
	matchedValue interface{}
	caseValues   []interface{}
	defaultCase  interface{}
	dict         *core.Thunk
}

// NewSwitch creates a switch expression.
func NewSwitch(m interface{}, cs []Case, d interface{}) Switch {
	if d == nil {
		panic("Default cases must be provided in switch expressions")
	}

	vs := make([]interface{}, 0, len(cs))

	for _, c := range cs {
		vs = append(vs, c.value)
	}

	return Switch{m, vs, d, compileCasesToDict(cs)}
}

func compileCasesToDict(cs []Case) *core.Thunk {
	ks := make([]core.Value, 0, len(cs))
	vs := make([]*core.Thunk, 0, len(cs))

	for i, c := range cs {
		ks = append(ks, c.pattern.Eval())
		vs = append(vs, core.NewNumber(float64(i)))
	}

	return core.Normal(core.NewDictionary(ks, vs).Eval())
}
