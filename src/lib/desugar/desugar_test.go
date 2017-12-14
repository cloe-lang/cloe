package desugar

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestDesugar(t *testing.T) {
	for _, ss := range [][]interface{}{
		{
			ast.NewLetVar(
				"foo",
				ast.NewPApp("prepend", []interface{}{"42", "foo"}, debug.NewGoInfo(0))),
		},
		{
			ast.NewLetVar(
				"foo",
				ast.NewAnonymousFunction(ast.NewSignature(nil, nil, "", nil, nil, ""), "123")),
		},
		{
			ast.NewLetFunction(
				"foo",
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				nil,
				ast.NewAnonymousFunction(ast.NewSignature(nil, nil, "", nil, nil, ""), "123"),
				debug.NewGoInfo(0)),
		},
		{
			ast.NewLetFunction(
				"foo",
				ast.NewSignature(nil, nil, "", nil, nil, ""),
				[]interface{}{
					ast.NewLetVar(
						"x",
						ast.NewAnonymousFunction(ast.NewSignature(nil, nil, "", nil, nil, ""), "123")),
				},
				"x",
				debug.NewGoInfo(0)),
		},
		{
			ast.NewLetFunction(
				"factorial",
				ast.NewSignature([]string{"n"}, nil, "", nil, nil, ""),
				nil,
				ast.NewMatch("n", []ast.MatchCase{
					ast.NewMatchCase("0", "1"),
					ast.NewMatchCase(
						"_",
						ast.NewPApp("*",
							[]interface{}{"n", ast.NewPApp(
								"factorial",
								[]interface{}{ast.NewPApp("-", []interface{}{"n", "1"}, debug.NewGoInfo(0))},
								debug.NewGoInfo(0))},
							debug.NewGoInfo(0))),
				}), debug.NewGoInfo(0)),
		},
		{
			ast.NewMutualRecursion([]ast.LetFunction{
				ast.NewLetFunction(
					"foo",
					ast.NewSignature(nil, nil, "", nil, nil, ""),
					nil,
					"bar",
					debug.NewGoInfo(0)),
				ast.NewLetFunction(
					"bar",
					ast.NewSignature(nil, nil, "", nil, nil, ""),
					nil,
					"foo",
					debug.NewGoInfo(0)),
			}, debug.NewGoInfo(0)),
		},
		{
			ast.NewLetVar(
				"v",
				ast.NewMatch("x", []ast.MatchCase{ast.NewMatchCase("y", "z")})),
		},
	} {
		t.Logf("%#v", ss)

		for _, s := range Desugar(ss) {
			ast.Convert(func(x interface{}) interface{} {
				switch x := x.(type) {
				case ast.AnonymousFunction:
					t.Fail()
				case ast.LetFunction:
					assert.Equal(t, 0, len(newNames(x.Name()).findInLetFunction(x)))

					for _, l := range x.Lets() {
						_, ok := l.(ast.LetFunction)
						assert.False(t, ok)
					}
				case ast.Match:
					t.Fail()
				case ast.MutualRecursion:
					t.Fail()
				}

				return nil
			}, s)
		}
	}
}
