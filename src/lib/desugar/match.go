package desugar

import (
	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func desugarMatchExpression(x interface{}) interface{} {
	switch x := x.(type) {
	case ast.App:
		return ast.NewApp(
			desugarMatchExpression(x.Function()),
			desugarMatchExpression(x.Arguments()).(ast.Arguments),
			x.DebugInfo())
	case ast.Arguments:
		ps := make([]ast.PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, desugarMatchExpression(p).(ast.PositionalArgument))
		}

		ks := make([]ast.KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, desugarMatchExpression(k).(ast.KeywordArgument))
		}

		ds := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, d := range x.ExpandedDicts() {
			ds = append(ds, desugarMatchExpression(d))
		}

		return ast.NewArguments(ps, ks, ds)
	case ast.KeywordArgument:
		return ast.NewKeywordArgument(x.Name(), desugarMatchExpression(x.Value()))
	case ast.LetFunction:
		ls := make([]interface{}, 0, len(x.Lets()))

		for _, l := range x.Lets() {
			ls = append(ls, desugarMatchExpression(l))
		}

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			ls,
			desugarMatchExpression(x.Body()),
			x.DebugInfo())
	case ast.LetVar:
		return ast.NewLetVar(x.Name(), desugarMatchExpression(x.Expr()))
	case ast.Match:
		return desugarMatchIntoApp(x)
	case ast.Output:
		return ast.NewOutput(desugarMatchExpression(x.Expr()), x.Expanded())
	case ast.PositionalArgument:
		return ast.NewPositionalArgument(desugarMatchExpression(x.Value()), x.Expanded())
	default:
		return x
	}
}

func desugarMatchIntoApp(m ast.Match) interface{} {
	return ast.NewApp(
		createMatchFunction(m.Cases()),
		ast.NewArguments(
			[]ast.PositionalArgument{ast.NewPositionalArgument(m.Value(), false)},
			nil,
			nil),
		debug.NewGoInfo(0))
}

func createMatchFunction(cs []ast.Case) interface{} {
	panic("Not implemented.")
}
