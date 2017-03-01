package compile

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/compile/env"
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/raviqqe/tisp/src/lib/ir"
	"log"
)

type compiler struct {
	env env.Environment
}

func newCompiler() compiler {
	return compiler{env: prelude.Child()}
}

func (c *compiler) compile(module []interface{}) []Output {
	outputs := make([]Output, 0, 8) // TODO: Best cap?

	for _, node := range module {
		switch x := node.(type) {
		case ast.LetConst:
			c.env.Set(x.Name(), c.exprToThunk(x.Expr()))
		case ast.LetFunction:
			c.env.Set(x.Name(), ir.CompileFunction(c.compileSignature(x.Signature()), []interface{}{}, c.exprToIR(x.Signature(), x.Body())))
		case ast.Output:
			outputs = append(outputs, NewOutput(c.exprToThunk(x.Expr()), x.Expanded()))
		default:
			panic(fmt.Sprint("Invalid instruction.", x))
		}
	}

	return outputs
}

func (c *compiler) exprToThunk(expr interface{}) *core.Thunk {
	switch x := expr.(type) {
	case string:
		return getOrError(c.env, x)
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

		return core.App(c.exprToThunk(x.Function()), core.NewArguments(ps, ks, ds))
	}

	panic(fmt.Sprintf("Invalid type as an expression. %#v", expr))
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

func (c *compiler) exprToIR(sig ast.Signature, expr interface{}) interface{} {
	switch x := expr.(type) {
	case string:
		i, err := sig.NameToIndex(x)

		if err == nil {
			return i
		}

		t, err := c.env.Get(x)

		if err == nil {
			return t
		}

		log.Fatalln(err.Error())
	case ast.App:
		args := x.Arguments()

		ps := make([]ir.PositionalArgument, len(args.Positionals()))
		for i, p := range args.Positionals() {
			ps[i] = ir.NewPositionalArgument(c.exprToIR(sig, p.Value()), p.Expanded())
		}

		ks := make([]ir.KeywordArgument, len(args.Keywords()))
		for i, k := range args.Keywords() {
			ks[i] = ir.NewKeywordArgument(k.Name(), c.exprToIR(sig, k.Value()))
		}

		ds := make([]interface{}, len(args.ExpandedDicts()))
		for i, d := range args.ExpandedDicts() {
			ds[i] = c.exprToIR(sig, d)
		}

		return ir.NewApp(c.exprToIR(sig, x.Function()), ir.NewArguments(ps, ks, ds))
	}

	panic(fmt.Sprint("Invalid type.", expr))
}

func getOrError(e env.Environment, s string) *core.Thunk {
	t, err := e.Get(s)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return t
}
