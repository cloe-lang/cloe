package ir

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/core"
)

// Switch represents a switch expression.
type Switch struct {
	value    interface{}
	patterns []core.Value
	values   []int
}

// NewSwitch creates a switch expression.
func NewSwitch(ps []core.Value, vs []int) Switch {
	if len(ps) != len(vs) {
		panic(fmt.Errorf(
			"A number of patterns (%d) doesn't match with a number of corresponding values (%d)",
			len(ps),
			len(vs)))
	} else if len(ps) == 0 {
		panic(fmt.Errorf("A number of patterns must be more than 0"))
	}

	return Switch{ps, vs}
}

func (s Switch) compileToDict() core.DictionaryType {
	ts := make([]*core.Thunk, 0, len(s.values))

	for _, v := range s.values {
		ts = append(ts, core.NewNumber(float64(v)))
	}

	return core.NewDictionary(s.patterns, ts).Eval().(core.DictionaryType)
}
