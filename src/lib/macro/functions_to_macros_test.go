package macro

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/std"
)

func TestConversionBetweenASTAndThunk(t *testing.T) {
	for _, a := range []interface{}{
		"foo",
		[]interface{}{"foo", "bar"},
		[]interface{}{"foo", []interface{}{"bar", "baz"}},
		[]interface{}{"foo", []interface{}{"bar", []interface{}{"baz", "blah"}}},
	} {
		assert.Equal(t, a, thunkToAST(astToThunk(a)))
	}
}

func TestFunctionsToMacros(t *testing.T) {
	for _, fs := range []map[string]*core.Thunk{
		map[string]*core.Thunk{"merge": core.Merge},
		map[string]*core.Thunk{"merge": core.Merge, "list": std.List},
	} {
		FunctionsToMacros(fs)
	}
}
