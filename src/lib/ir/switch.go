package ir

import (
	"github.com/cloe-lang/cloe/src/lib/core"
)

// Switch represents a switch expression.
type Switch struct {
	matchedValue any
	caseValues   []any
	defaultCase  any
	dict         core.Value
}

// NewSwitch creates a switch expression.
func NewSwitch(m any, cs []Case, d any) Switch {
	if d == nil {
		panic("Default cases must be provided in switch expressions")
	}

	vs := make([]any, 0, len(cs))

	for _, c := range cs {
		vs = append(vs, c.value)
	}

	return Switch{m, vs, d, compileCasesToDictionary(cs)}
}

func compileCasesToDictionary(cs []Case) core.Value {
	kvs := make([]core.KeyValue, 0, len(cs))

	for i, c := range cs {
		kvs = append(kvs, core.KeyValue{Key: c.pattern, Value: core.NewNumber(float64(i))})
	}

	return core.EvalPure(core.NewDictionary(kvs))
}
