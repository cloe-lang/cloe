package desugar

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestNamesFindInLetVar(t *testing.T) {
	n := "x"
	assert.True(t, newNames(n).findInLetVar(ast.NewLetVar(n, n)).include(n))
}

func TestNamesFindInDefFunction(t *testing.T) {
	n := "x"

	for _, test := range []struct {
		letFunc ast.DefFunction
		answer  bool
	}{
		{
			ast.NewDefFunction(
				n,
				ast.NewSignature(nil, "", nil, ""),
				nil,
				n,
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewDefFunction(
				n,
				ast.NewSignature(nil, "", nil, ""),
				[]interface{}{ast.NewLetVar(n, "y")},
				n,
				debug.NewGoInfo(0)),
			false,
		},
		{
			ast.NewDefFunction(
				n,
				ast.NewSignature(nil, "", nil, ""),
				[]interface{}{ast.NewLetVar(n, "x")},
				"42",
				debug.NewGoInfo(0)),
			true,
		},
	} {
		assert.Equal(t, test.answer, newNames(n).findInDefFunction(test.letFunc).include(n))
	}
}

func TestNamesFindInDefFunctionPanic(t *testing.T) {
	assert.Panics(t, func() {
		newNames().findInDefFunction(ast.NewDefFunction(
			"func",
			ast.NewSignature(nil, "", nil, ""),
			[]interface{}{nil},
			"x",
			debug.NewGoInfo(0)))
	})
}

func TestNamesFindInExpression(t *testing.T) {
	n := "x"

	for _, c := range []struct {
		expression interface{}
		answer     bool
	}{
		{
			"x",
			true,
		},
		{
			ast.NewPApp("f", []interface{}{"x"}, debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewApp(
				"f",
				ast.NewArguments(nil, []ast.KeywordArgument{ast.NewKeywordArgument("foo", "x")}),
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewApp(
				"f",
				ast.NewArguments(nil, []ast.KeywordArgument{ast.NewKeywordArgument("", "x")}),
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewAnonymousFunction(ast.NewSignature([]string{"x"}, "", nil, ""), nil, "x"),
			false,
		},
	} {
		assert.Equal(t, c.answer, newNames(n).findInExpression(c.expression).include(n))
	}
}

func TestNamesFindInExpressionPanic(t *testing.T) {
	assert.Panics(t, func() {
		newNames().findInExpression(nil)
	})
}
