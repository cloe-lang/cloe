package desugar

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

var letFooFunction = ast.NewLetFunction(
	"foo",
	ast.NewSignature(nil, nil, "", nil, nil, ""),
	nil,
	"nil",
	debug.NewGoInfo(0))

func TestDesugarMutualRecursion(t *testing.T) {
	for _, mr := range []ast.MutualRecursion{
		ast.NewMutualRecursion(
			[]ast.LetFunction{
				ast.NewLetFunction(
					"foo",
					ast.NewSignature(nil, nil, "", nil, nil, ""),
					nil,
					"nil",
					debug.NewGoInfo(0)),
				ast.NewLetFunction(
					"bar",
					ast.NewSignature(nil, nil, "", nil, nil, ""),
					nil,
					"nil",
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
		ast.NewMutualRecursion(
			[]ast.LetFunction{
				ast.NewLetFunction(
					"foo",
					ast.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
					nil,
					ast.NewPApp("bar", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
				ast.NewLetFunction(
					"bar",
					ast.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
					nil,
					ast.NewPApp("foo", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
		ast.NewMutualRecursion(
			[]ast.LetFunction{
				ast.NewLetFunction(
					"foo",
					ast.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
					[]interface{}{
						ast.NewLetFunction(
							"f",
							ast.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
							nil,
							ast.NewPApp("bar", []interface{}{"x"}, debug.NewGoInfo(0)),
							debug.NewGoInfo(0)),
						ast.NewLetVar("g", "f"),
					},
					ast.NewPApp("g", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
				ast.NewLetFunction(
					"bar",
					ast.NewSignature([]string{"x"}, nil, "", nil, nil, ""),
					nil,
					ast.NewPApp("foo", []interface{}{"x"}, debug.NewGoInfo(0)),
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
	} {
		desugarMutualRecursion(mr)
	}
}

func TestDesugarMutualRecursionWithOneFunction(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	desugarMutualRecursion(ast.NewMutualRecursion(
		[]ast.LetFunction{letFooFunction},
		debug.NewGoInfo(0)))
}

func TestDesugarMutualRecursionWithFunctionsOfSameName(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	desugarMutualRecursion(ast.NewMutualRecursion(
		[]ast.LetFunction{letFooFunction, letFooFunction},
		debug.NewGoInfo(0)))
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
				case ast.LetFunction:
					name = l.Name()
				case ast.LetVar:
					name = l.Name()
				default:
					panic("Unreachable")
				}

				assert.Equal(t, name, s)
			}
		}
	}
}
