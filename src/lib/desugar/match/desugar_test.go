package match

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func TestDesugar(t *testing.T) {
	for _, s := range []interface{}{
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
		ast.NewMutualRecursion([]ast.LetFunction{
			ast.NewLetFunction(
				"even?",
				ast.NewSignature([]string{"n"}, nil, "", nil, nil, ""),
				nil,
				ast.NewMatch("n", []ast.MatchCase{
					ast.NewMatchCase("0", "true"),
					ast.NewMatchCase(
						"_",
						ast.NewPApp(
							"odd?",
							[]interface{}{ast.NewPApp("-", []interface{}{"n", "1"}, debug.NewGoInfo(0))},
							debug.NewGoInfo(0))),
				}), debug.NewGoInfo(0)),
			ast.NewLetFunction(
				"odd?",
				ast.NewSignature([]string{"n"}, nil, "", nil, nil, ""),
				nil,
				ast.NewMatch("n", []ast.MatchCase{
					ast.NewMatchCase("0", "true"),
					ast.NewMatchCase(
						"_",
						ast.NewPApp(
							"even?",
							[]interface{}{ast.NewPApp("-", []interface{}{"n", "1"}, debug.NewGoInfo(0))},
							debug.NewGoInfo(0))),
				}), debug.NewGoInfo(0)),
		}, debug.NewGoInfo(0)),
	} {
		for _, s := range Desugar(s) {
			t.Logf("%#v", s)

			ast.Convert(func(x interface{}) interface{} {
				if _, ok := x.(ast.Match); ok {
					t.Fail()
				}

				return nil
			}, s)
		}
	}
}
