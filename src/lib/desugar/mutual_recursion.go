package desugar

import (
	"fmt"

	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/debug"
	"github.com/raviqqe/tisp/src/lib/gensym"
	"github.com/raviqqe/tisp/src/lib/util"
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
	olds := mr.LetFunctions()
	nonRecursives := make([]interface{}, 0, len(olds))

	for _, old := range olds {
		fsArg := gensym.GenSym("mr", "functions", "argument")
		nameToIndex := indexLetFunctions(olds...)

		nonRecursives = append(
			nonRecursives,
			ast.NewLetFunction(
				gensym.GenSym("nonRecursive", old.Name()),
				prependPosReqsToSig(old.Signature(), []string{fsArg}),
				replaceNames(fsArg, nameToIndex, old.Lets(), mr.DebugInfo()).([]interface{}),
				replaceNames(fsArg, deleteNamesDefinedByLets(nameToIndex, old.Lets()), old.Body(), mr.DebugInfo()),
				old.DebugInfo()))
	}

	mrFunctionList := gensym.GenSym("ys", "mr", "functions")
	news := make([]interface{}, 0, len(olds))

	for i, old := range olds {
		news = append(
			news,
			ast.NewLetVar(
				old.Name(),
				ast.NewPApp(mrFunctionList, []interface{}{fmt.Sprint(i)}, old.DebugInfo())))
	}

	return append(
		nonRecursives,
		append(
			[]interface{}{ast.NewLetVar(
				mrFunctionList,
				ast.NewPApp("$ys", util.StringsToAnys(letStatementsToNames(nonRecursives)), mr.DebugInfo()))},
			news...)...)
}

func indexLetFunctions(fs ...ast.LetFunction) map[string]int {
	nameToIndex := make(map[string]int)

	for i, f := range fs {
		nameToIndex[f.Name()] = i
	}

	if len(nameToIndex) != len(fs) {
		util.Fail("Duplicate names were found among mutually-recursive functions.")
	}

	return nameToIndex
}

func replaceNames(functionList string, nameToIndex map[string]int, x interface{}, di debug.Info) interface{} {
	replaceWithNameToIndex := func(nameToIndex map[string]int) func(x interface{}) interface{} {
		return func(x interface{}) interface{} {
			return replaceNames(functionList, nameToIndex, x, di)
		}
	}

	replace := replaceWithNameToIndex(nameToIndex)

	switch x := x.(type) {
	case []interface{}:
		ys := make([]interface{}, 0, len(x))

		for _, x := range x {
			ys = append(ys, replace(x))
		}

		return ys
	case ast.LetFunction:
		nameToIndex := copyNameToIndex(nameToIndex)

		delete(nameToIndex, x.Name())
		for _, n := range signatureToNames(x.Signature()).slice() {
			delete(nameToIndex, n)
		}

		return ast.NewLetFunction(
			x.Name(),
			x.Signature(),
			replaceWithNameToIndex(nameToIndex)(x.Lets()).([]interface{}),
			replaceWithNameToIndex(deleteNamesDefinedByLets(nameToIndex, x.Lets()))(x.Body()),
			x.DebugInfo())
	case ast.LetVar:
		nameToIndex := copyNameToIndex(nameToIndex)
		delete(nameToIndex, x.Name())
		return ast.NewLetVar(x.Name(), replaceWithNameToIndex(nameToIndex)(x.Expr()))
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
		if i, ok := nameToIndex[x]; ok {
			return ast.NewPApp(functionList, []interface{}{fmt.Sprint(i)}, di)
		}

		return x
	}

	panic(x)
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
	names := make([]string, 0, len(ls))

	for _, l := range ls {
		switch l := l.(type) {
		case ast.LetFunction:
			names = append(names, l.Name())
		case ast.LetVar:
			names = append(names, l.Name())
		default:
			panic("Unreachable")
		}
	}

	return names
}
