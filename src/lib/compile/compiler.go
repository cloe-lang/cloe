package compile

import (
	"fmt"
	"path"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/ir"
)

type compiler struct {
	env environment
}

func newCompiler(e environment) compiler {
	return compiler{e}
}

func (c *compiler) compile(module []interface{}) []Output {
	outputs := make([]Output, 0)

	for _, s := range module {
		switch x := s.(type) {
		case ast.LetVar:
			c.env.set(x.Name(), c.exprToThunk(x.Expr()))
		case ast.LetFunction:
			sig := x.Signature()
			ls := x.Lets()

			vars := make([]interface{}, 0, len(ls))
			varToIndex := sig.NameToIndex()

			for _, l := range ls {
				v := l.(ast.LetVar)
				vars = append(vars, c.exprToIR(varToIndex, v.Expr()))
				varToIndex[v.Name()] = len(varToIndex)
			}

			c.env.set(
				x.Name(),
				ir.CompileFunction(
					c.compileSignature(sig),
					vars,
					c.exprToIR(varToIndex, x.Body())))
		case ast.Output:
			outputs = append(outputs, NewOutput(c.exprToThunk(x.Expr()), x.Expanded()))
		case ast.Import:
			for k, v := range SubModule(x.Path() + ".tisp") {
				c.env.set(path.Base(x.Path())+"."+k, v)
			}
		default:
			panic(fmt.Errorf("Invalid type: %#v", x))
		}
	}

	return outputs
}

func (c *compiler) exprToThunk(expr interface{}) *core.Thunk {
	return core.PApp(ir.CompileFunction(
		core.NewSignature(nil, nil, "", nil, nil, ""),
		nil,
		c.exprToIR(nil, expr)))
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

func (c *compiler) exprToIR(varToIndex map[string]int, expr interface{}) interface{} {
	switch x := expr.(type) {
	case string:
		if i, ok := varToIndex[x]; ok {
			return i
		}

		return c.env.get(x)
	case ast.App:
		args := x.Arguments()

		ps := make([]ir.PositionalArgument, 0, len(args.Positionals()))
		for _, p := range args.Positionals() {
			ps = append(ps, ir.NewPositionalArgument(c.exprToIR(varToIndex, p.Value()), p.Expanded()))
		}

		ks := make([]ir.KeywordArgument, 0, len(args.Keywords()))
		for _, k := range args.Keywords() {
			ks = append(ks, ir.NewKeywordArgument(k.Name(), c.exprToIR(varToIndex, k.Value())))
		}

		ds := make([]interface{}, 0, len(args.ExpandedDicts()))
		for _, d := range args.ExpandedDicts() {
			ds = append(ds, c.exprToIR(varToIndex, d))
		}

		return ir.NewApp(
			c.exprToIR(varToIndex, x.Function()),
			ir.NewArguments(ps, ks, ds),
			x.DebugInfo())
	case ast.Switch:
		cs := make([]ir.Case, 0, len(x.Cases()))

		for _, k := range x.Cases() {
			cs = append(cs, ir.NewCase(
				c.env.get(k.Pattern()),
				c.exprToIR(varToIndex, k.Value())))
		}

		d := interface{}(nil)

		if x.DefaultCase() != nil {
			d = c.exprToIR(varToIndex, x.DefaultCase())
		}

		return ir.NewSwitch(c.exprToIR(varToIndex, x.Value()), cs, d)
	}

	panic(fmt.Errorf("Invalid type: %#v", expr))
}
