package compile

import (
	"../ast"
	"../ir"
	"../vm"
	"./env"
	"fmt"
	"log"
)

type compiler struct {
	env env.Environment
}

func newCompiler() compiler {
	return compiler{env: prelude.Child()}
}

func (c *compiler) compile(module []interface{}) []*vm.Thunk {
	outputs := make([]*vm.Thunk, 0, 8) // TODO: Use ListType.

	for _, node := range module {
		switch x := node.(type) {
		case ast.LetConst:
			c.env.Set(x.Name(), c.exprToThunk(x.Expr()))
		case ast.LetFunction:
			c.env.Set(x.Name(), ir.CompileFunction(c.compileSignature(x.Signature()), c.exprToIR(x.Signature(), x.Body())))
		case ast.Output:
			outputs = append(outputs, c.exprToThunk(x.Expr())) // TODO: Expand expanded lists.
		default:
			panic(fmt.Sprint("Invalid instruction.", x))
		}
	}

	return outputs
}

func (c *compiler) exprToThunk(expr interface{}) *vm.Thunk {
	switch x := expr.(type) {
	case string:
		return getOrError(c.env, x)
	case ast.App:
		args := x.Arguments()

		ps := make([]vm.PositionalArgument, len(args.Positionals()))
		for i, p := range args.Positionals() {
			ps[i] = vm.NewPositionalArgument(c.exprToThunk(p.Value()), p.Expanded())
		}

		ks := make([]vm.KeywordArgument, len(args.Keywords()))
		for i, k := range args.Keywords() {
			ks[i] = vm.NewKeywordArgument(k.Name(), c.exprToThunk(k.Value()))
		}

		ds := make([]*vm.Thunk, len(args.ExpandedDicts()))
		for i, d := range args.ExpandedDicts() {
			ds[i] = c.exprToThunk(d)
		}

		return vm.App(c.exprToThunk(x.Function()), vm.NewArguments(ps, ks, ds))
	}

	panic(fmt.Sprintf("Invalid type as an expression. %#v", expr))
}

func (c *compiler) compileSignature(sig ast.Signature) vm.Signature {
	return vm.NewSignature(
		sig.PosReqs(), c.compileOptionalArguments(sig.PosOpts()), sig.PosRest(),
		sig.KeyReqs(), c.compileOptionalArguments(sig.KeyOpts()), sig.KeyRest(),
	)
}

func (c *compiler) compileOptionalArguments(opts []ast.OptionalArgument) []vm.OptionalArgument {
	vmOpts := make([]vm.OptionalArgument, len(opts))

	for i, opt := range opts {
		vmOpts[i] = vm.NewOptionalArgument(opt.Name(), c.exprToThunk(opt.DefaultValue()))
	}

	return vmOpts
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

func getOrError(e env.Environment, s string) *vm.Thunk {
	t, err := e.Get(s)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return t
}
