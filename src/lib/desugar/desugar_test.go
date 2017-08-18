package desugar

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
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
	} {
		t.Logf("%#v", ss)

		for _, s := range Desugar(ss) {
			ast.Convert(func(x interface{}) interface{} {
				switch x := x.(type) {
				case ast.AnonymousFunction:
					t.Fail()
				case ast.LetFunction:
					if len(newNames(x.Name()).findInFunction(x)) != 0 {
						t.Fail()
					}

					for _, l := range x.Lets() {
						if _, ok := l.(ast.LetFunction); ok {
							t.Fail()
						}
					}
				case ast.LetVar:
					if len(newNames(x.Name()).find(x.Expr())) != 0 {
						t.Fail()
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
