package desugar

import (
	"fmt"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/gensym"
)

func flattenStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.DefFunction:
		return flattenDefFunction(s)
	default:
		return []interface{}{s}
	}
}

func flattenStatements(old []interface{}) []interface{} {
	new := make([]interface{}, 0)

	for _, s := range old {
		new = append(new, flattenStatement(s)...)
	}

	return new
}

func flattenDefFunction(f ast.DefFunction) []interface{} {
	f = flattenInnerStatements(f)

	ss := make([]interface{}, 0)
	ls := make([]interface{}, 0)
	ns := signatureToNames(f.Signature())

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetVar:
			ls = append(ls, l)
			ns.add(l.Name())
		case ast.DefFunction:
			args := ns.findInDefFunction(l).slice()
			n := gensym.GenSym()

			ss = append(ss, letFlattenedFunction(l, n, args))
			ls = append(ls, letClosure(l, n, args))

			ns.add(l.Name())
		default:
			panic(fmt.Errorf("Invalid value: %#v", l))
		}
	}

	return append(ss, ast.NewDefFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo()))
}

func letFlattenedFunction(f ast.DefFunction, n string, args []string) ast.DefFunction {
	return ast.NewDefFunction(
		n,
		prependPositionalsToSig(args, f.Signature()),
		f.Lets(),
		f.Body(),
		f.DebugInfo())
}

func letClosure(f ast.DefFunction, n string, args []string) ast.LetVar {
	return ast.NewLetVar(
		f.Name(),
		ast.NewApp(
			"$partial",
			ast.NewArguments(namesToPosArgs(append([]string{n}, args...)), nil),
			f.DebugInfo()))
}

func namesToPosArgs(ns []string) []ast.PositionalArgument {
	ps := make([]ast.PositionalArgument, 0, len(ns))

	for _, n := range ns {
		ps = append(ps, ast.NewPositionalArgument(n, false))
	}

	return ps
}

func flattenInnerStatements(f ast.DefFunction) ast.DefFunction {
	return ast.NewDefFunction(
		f.Name(),
		f.Signature(),
		flattenStatements(f.Lets()),
		f.Body(),
		f.DebugInfo())
}
