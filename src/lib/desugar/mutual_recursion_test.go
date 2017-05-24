package desugar

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/stretchr/testify/assert"
)

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
