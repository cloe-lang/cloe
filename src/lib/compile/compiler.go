package compile

import (
	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/raviqqe/tisp/src/lib/ir"
	"github.com/raviqqe/tisp/src/lib/util"
)

type compiler struct {
	env environment
}

func newCompiler() compiler {
	return compiler{env: prelude()}
}

func (c *compiler) compile(module []interface{}) []Output {
	outputs := make([]Output, 0, 8) // TODO: Best cap?

	for _, node := range module {
		switch x := node.(type) {
		case ast.LetConst:
			c.env.set(x.Name(), c.exprToThunk(x.Expr()))
		case ast.LetFunction:
			sig := x.Signature()
			ls := x.Lets()

			vars := make([]interface{}, len(ls))
			varIndices := map[string]int{}

			for i, l := range ls {
				cst := l.(ast.LetConst)
				vars[i] = c.exprToIR(sig, varIndices, cst.Expr())
				varIndices[cst.Name()] = sig.Arity() + i
			}

			c.env.set(
				x.Name(),
				ir.CompileFunction(
					c.compileSignature(sig),
					vars,
					c.exprToIR(sig, varIndices, x.Body())))
		case ast.Output:
			outputs = append(outputs, NewOutput(c.exprToThunk(x.Expr()), x.Expanded()))
		default:
			util.Fail("Invalid instruction.", x)
		}
	}

	return outputs
}

func (c *compiler) exprToThunk(expr interface{}) *core.Thunk {
	switch x := expr.(type) {
	case string:
		return c.env.get(x)
	case ast.App:
		args := x.Arguments()

		ps := make([]core.PositionalArgument, len(args.Positionals()))
		for i, p := range args.Positionals() {
			ps[i] = core.NewPositionalArgument(c.exprToThunk(p.Value()), p.Expanded())
		}

		ks := make([]core.KeywordArgument, len(args.Keywords()))
		for i, k := range args.Keywords() {
			ks[i] = core.NewKeywordArgument(k.Name(), c.exprToThunk(k.Value()))
		}

		ds := make([]*core.Thunk, len(args.ExpandedDicts()))
		for i, d := range args.ExpandedDicts() {
			ds[i] = c.exprToThunk(d)
		}

		return core.AppWithInfo(
			c.exprToThunk(x.Function()),
			core.NewArguments(ps, ks, ds),
			x.DebugInfo())
	}

	util.Fail("Invalid type as an expression. %#v", expr)
	panic("Unreachable")
}

func (c *compiler) compileSignature(sig ast.Signature) core.Signature {
	return core.NewSignature(
		sig.PosReqs(), c.compileOptionalArguments(sig.PosOpts()), sig.PosRest(),
		sig.KeyReqs(), c.compileOptionalArguments(sig.KeyOpts()), sig.KeyRest(),
	)
}

func (c *compiler) compileOptionalArguments(opts []ast.OptionalArgument) []core.OptionalArgument {
	coreOpts := make([]core.OptionalArgument, len(opts))

	for i, opt := range opts {
		coreOpts[i] = core.NewOptionalArgument(opt.Name(), c.exprToThunk(opt.DefaultValue()))
	}

	return coreOpts
}

func (c *compiler) exprToIR(sig ast.Signature, vars map[string]int, expr interface{}) interface{} {
	switch x := expr.(type) {
	case string:
		if i, ok := vars[x]; ok {
			return i
		}

		i, err := sig.NameToIndex(x)

		if err == nil {
			return i
		}

		return c.env.get(x)
	case ast.App:
		args := x.Arguments()

		ps := make([]ir.PositionalArgument, len(args.Positionals()))
		for i, p := range args.Positionals() {
			ps[i] = ir.NewPositionalArgument(c.exprToIR(sig, vars, p.Value()), p.Expanded())
		}

		ks := make([]ir.KeywordArgument, len(args.Keywords()))
		for i, k := range args.Keywords() {
			ks[i] = ir.NewKeywordArgument(k.Name(), c.exprToIR(sig, vars, k.Value()))
		}

		ds := make([]interface{}, len(args.ExpandedDicts()))
		for i, d := range args.ExpandedDicts() {
			ds[i] = c.exprToIR(sig, vars, d)
		}

		return ir.NewApp(
			c.exprToIR(sig, vars, x.Function()),
			ir.NewArguments(ps, ks, ds),
			x.DebugInfo())
	}

	util.Fail("Invalid type.", expr)
	panic("Unreachable")
}
