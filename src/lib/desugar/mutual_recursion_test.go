package desugar

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

var letFooFunction = ast.NewDefFunction(
	"foo",
	ast.NewSignature(nil, "", nil, ""),
	nil,
	"nil",
	debug.NewGoInfo(0))

func TestDesugarMutualRecursion(t *testing.T) {
	for _, mr := range []ast.MutualRecursion{
		ast.NewMutualRecursion(
			[]ast.DefFunction{
				ast.NewDefFunction(
					"foo",
					ast.NewSignature(nil, "", nil, ""),
					nil,
					"nil",
					debug.NewGoInfo(0)),
				ast.NewDefFunction(
					"bar",
					ast.NewSignature(nil, "", nil, ""),
					nil,
					"nil",
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
		ast.NewMutualRecursion(
			[]ast.DefFunction{
				ast.NewDefFunction(
					"foo",
					ast.NewSignature([]string{"x"}, "", nil, ""),
					nil,
					ast.NewPApp("bar", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
				ast.NewDefFunction(
					"bar",
					ast.NewSignature([]string{"x"}, "", nil, ""),
					nil,
					ast.NewPApp("foo", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
		ast.NewMutualRecursion(
			[]ast.DefFunction{
				ast.NewDefFunction(
					"foo",
					ast.NewSignature([]string{"x"}, "", nil, ""),
					[]interface{}{
						ast.NewDefFunction(
							"f",
							ast.NewSignature([]string{"x"}, "", nil, ""),
							nil,
							ast.NewPApp("bar", []interface{}{"x"}, debug.NewGoInfo(0)),
							debug.NewGoInfo(0)),
						ast.NewLetVar("g", "f"),
					},
					ast.NewPApp("g", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
				ast.NewDefFunction(
					"bar",
					ast.NewSignature([]string{"x"}, "", nil, ""),
					nil,
					ast.NewPApp("foo", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
	} {
		desugarMutualRecursion(mr)
	}
}

func TestDesugarMutualRecursionWithOneFunction(t *testing.T) {
	assert.Panics(t, func() {
		desugarMutualRecursion(ast.NewMutualRecursion(
			[]ast.DefFunction{letFooFunction},
			debug.NewGoInfo(0)))
	})
}

func TestDesugarMutualRecursionWithFunctionsOfSameName(t *testing.T) {
	assert.Panics(t, func() {
		desugarMutualRecursion(ast.NewMutualRecursion(
			[]ast.DefFunction{letFooFunction, letFooFunction},
			debug.NewGoInfo(0)))
	})
}

func TestLetStatementsToNames(t *testing.T) {
	for i := 0; i < 100; i++ {
		for _, ls := range [][]interface{}{
			{},
			{ast.NewLetVar("foo", "bar")},
			{ast.NewLetVar("foo", "nil"), ast.NewLetVar("bar", "nil")},
			{ast.NewLetVar("foo", "nil"), ast.NewLetVar("bar", "nil"), ast.NewLetVar("baz", "nil")},
			{ast.NewLetVar("foo0", "nil"), ast.NewLetVar("foo1", "nil"), ast.NewLetVar("foo3", "nil"),
				ast.NewLetVar("foo4", "nil"), ast.NewLetVar("foo5", "nil"), ast.NewLetVar("foo6", "nil")},
		} {
			for i, s := range letStatementsToNames(ls) {
				var name string

				switch l := ls[i].(type) {
				case ast.DefFunction:
					name = l.Name()
				case ast.LetVar:
					name = l.Name()
				default:
					panic("unreachable")
				}

				assert.Equal(t, name, s)
			}
		}
	}
}
