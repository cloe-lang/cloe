package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/consts"
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
				ast.NewMatchCase("_", papp("*", "n", papp("factorial", papp("-", "n", "1")))),
			}), debug.NewGoInfo(0)),
		ast.NewMutualRecursion([]ast.LetFunction{
			ast.NewLetFunction(
				"even?",
				ast.NewSignature([]string{"n"}, nil, "", nil, nil, ""),
				nil,
				ast.NewMatch("n", []ast.MatchCase{
					ast.NewMatchCase("0", "true"),
					ast.NewMatchCase("_", papp("odd?", papp("-", "n", "1"))),
				}), debug.NewGoInfo(0)),
			ast.NewLetFunction(
				"odd?",
				ast.NewSignature([]string{"n"}, nil, "", nil, nil, ""),
				nil,
				ast.NewMatch("n", []ast.MatchCase{
					ast.NewMatchCase("0", "true"),
					ast.NewMatchCase("_", papp("even?", papp("-", "n", "1"))),
				}), debug.NewGoInfo(0)),
		}, debug.NewGoInfo(0)),
		ast.NewLetVar("x", ast.NewMatch("nil", []ast.MatchCase{
			ast.NewMatchCase(papp(consts.Names.ListFunction, "1", "x"), "x"),
			ast.NewMatchCase(papp(consts.Names.DictionaryFunction, "1", "x", `"foo"`, "true"), "x"),
		})),
	} {
		for _, s := range Desugar(s) {
			t.Logf("%#v", s)

			ast.Convert(func(x interface{}) interface{} {
				_, ok := x.(ast.Match)
				assert.False(t, ok)

				return nil
			}, s)
		}
	}
}

func papp(xs ...interface{}) ast.App {
	return ast.NewPApp(xs[0], xs[1:], debug.NewGoInfo(1))
}
