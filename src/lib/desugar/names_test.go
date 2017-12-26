package desugar

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestNamesFindInLetVar(t *testing.T) {
	n := "x"
	assert.True(t, newNames(n).findInLetVar(ast.NewLetVar(n, n)).include(n))
}

func TestNamesFindInLetFunction(t *testing.T) {
	n := "x"

	for _, test := range []struct {
		letFunc ast.LetFunction
		answer  bool
	}{
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				nil,
				n,
				debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewLetFunction(
				n,
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				[]interface{}{ast.NewLetVar(n, "y")},
				n,
				debug.NewGoInfo(0)),
			false,
		},
	} {
		assert.Equal(t, test.answer, newNames(n).findInLetFunction(test.letFunc).include(n))
	}
}

func TestNamesFindInLetFunctionPanic(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	newNames().findInLetFunction(ast.NewLetFunction(
		"func",
		ast.NewSignature(nil, nil, "", nil, nil, ""),
		[]interface{}{nil},
		"x",
		debug.NewGoInfo(0)))
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
			ast.NewApp("f", ast.NewArguments(nil, []ast.KeywordArgument{
				ast.NewKeywordArgument("foo", "x"),
			}, nil), debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewApp("f", ast.NewArguments(nil, nil, []interface{}{"x"}), debug.NewGoInfo(0)),
			true,
		},
		{
			ast.NewAnonymousFunction(ast.NewSignature([]string{"x"}, nil, "", nil, nil, ""), "x"),
			false,
		},
	} {
		assert.Equal(t, c.answer, newNames(n).findInExpression(c.expression).include(n))
	}
}

func TestNamesFindInExpressionPanic(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	newNames().findInExpression(nil)
}
