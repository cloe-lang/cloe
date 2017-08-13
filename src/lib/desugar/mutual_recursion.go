package desugar

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func desugarMutualRecursionStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.MutualRecursion:
		return desugarMutualRecursion(s)
	default:
		return []interface{}{s}
	}
}

func desugarMutualRecursion(mr ast.MutualRecursion) []interface{} {
	fs := mr.LetFunctions()
	unrecs := make([]interface{}, 0, len(fs))

	for _, f := range fs {
		arg := gensym.GenSym("mr", "functions", "argument")
		n2i := indexLetFunctions(fs...)

		unrecs = append(
			unrecs,
			ast.NewLetFunction(
				gensym.GenSym("mr", "unrec", f.Name()),
				prependPosReqsToSig([]string{arg}, f.Signature()),
				replaceNames(arg, n2i, f.Lets(), mr.DebugInfo()).([]interface{}),
				replaceNames(arg, deleteNamesDefinedByLets(n2i, f.Lets()), f.Body(), mr.DebugInfo()),
				f.DebugInfo()))
	}

	recsList := gensym.GenSym("ys", "mr", "functions")
	recs := make([]interface{}, 0, len(fs))

	for i, f := range fs {
		recs = append(
			recs,
			ast.NewLetVar(
				f.Name(),
				ast.NewPApp(recsList, []interface{}{fmt.Sprint(i)}, f.DebugInfo())))
	}

	return append(
		unrecs,
		append(
			[]interface{}{ast.NewLetVar(
				recsList,
				ast.NewPApp("$ys", stringsToAnys(letStatementsToNames(unrecs)), mr.DebugInfo()))},
			recs...)...)
}

func indexLetFunctions(fs ...ast.LetFunction) map[string]int {
	n2i := make(map[string]int)

	for i, f := range fs {
		n2i[f.Name()] = i
	}

	if len(n2i) != len(fs) {
		panic(fmt.Errorf("Duplicate names were found among mutually-recursive functions"))
	}

	return n2i
}

func replaceNames(funcList string, n2i map[string]int, x interface{}, di debug.Info) interface{} {
	replaceWithNameToIndex := func(n2i map[string]int) func(x interface{}) interface{} {
		return func(x interface{}) interface{} {
			return replaceNames(funcList, n2i, x, di)
		}
	}

	replace := replaceWithNameToIndex(n2i)

	switch x := x.(type) {
	case []interface{}:
		ys := make([]interface{}, 0, len(x))

		for _, x := range x {
			ys = append(ys, replace(x))
		}

		return ys
	case ast.LetFunction:
		n2i := copyNameToIndex(n2i)

		delete(n2i, x.Name())
		for n := range signatureToNames(x.Signature()) {
			delete(n2i, n)
		}

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			replaceWithNameToIndex(n2i)(x.Lets()).([]interface{}),
			replaceWithNameToIndex(deleteNamesDefinedByLets(n2i, x.Lets()))(x.Body()),
			x.DebugInfo())
	case ast.LetVar:
		n2i := copyNameToIndex(n2i)
		delete(n2i, x.Name())
		return ast.NewLetVar(x.Name(), replaceWithNameToIndex(n2i)(x.Expr()))
	case ast.App:
		return ast.NewApp(replace(x.Function()), replace(x.Arguments()).(ast.Arguments), x.DebugInfo())
	case ast.Arguments:
		ps := make([]ast.PositionalArgument, 0, len(x.Positionals()))

		for _, p := range x.Positionals() {
			ps = append(ps, ast.NewPositionalArgument(replace(p.Value()), p.Expanded()))
		}

		ks := make([]ast.KeywordArgument, 0, len(x.Keywords()))

		for _, k := range x.Keywords() {
			ks = append(ks, ast.NewKeywordArgument(k.Name(), replace(k.Value())))
		}

		ds := make([]interface{}, 0, len(x.ExpandedDicts()))

		for _, d := range x.ExpandedDicts() {
			ds = append(ds, replace(d))
		}

		return ast.NewArguments(ps, ks, ds)
	case string:
		if i, ok := n2i[x]; ok {
			return ast.NewPApp(funcList, []interface{}{fmt.Sprint(i)}, di)
		}

		return x
	}

	panic(fmt.Errorf("Invalid value: %#v", x))
}

func copyNameToIndex(ni map[string]int) map[string]int {
	new := make(map[string]int)

	for k, v := range ni {
		new[k] = v
	}

	return new
}

func deleteNamesDefinedByLets(ni map[string]int, ls []interface{}) map[string]int {
	ni = copyNameToIndex(ni)

	for _, n := range letStatementsToNames(ls) {
		delete(ni, n)
	}

	return ni
}

func letStatementsToNames(ls []interface{}) []string {
	ns := make([]string, 0, len(ls))

	for _, l := range ls {
		switch l := l.(type) {
		case ast.LetFunction:
			ns = append(ns, l.Name())
		case ast.LetVar:
			ns = append(ns, l.Name())
		default:
			panic("Unreachable")
		}
	}

	return ns
}

func stringsToAnys(ss []string) []interface{} {
	xs := make([]interface{}, 0, len(ss))

	for _, s := range ss {
		xs = append(xs, s)
	}

	return xs
}
