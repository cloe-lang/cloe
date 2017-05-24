package desugar

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/gensym"
)

func flattenStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.LetFunction:
		return flattenLetFunction(s)
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

func flattenLetFunction(f ast.LetFunction) []interface{} {
	f = flattenInnerStatements(f)

	ss := make([]interface{}, 0)
	ls := make([]interface{}, 0)
	names := signatureToNames(f.Signature())

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetVar:
			ls = append(ls, l)
			names.add(l.Name())
		case ast.LetFunction:
			usedNames := names.find(l.Body())
			for _, l := range l.Lets() {
				usedNames.merge(names.find(l.(ast.LetVar)))
			}
			usedNames.subtract(signatureToNames(l.Signature()))

			flattened := gensym.GenSym(f.Name(), l.Name())

			ss = append(ss, ast.NewLetFunction(
				flattened,
				prependPosReqsToSig(l.Signature(), usedNames.slice()),
				l.Lets(),
				l.Body(),
				l.DebugInfo()))

			ls = append(ls, ast.NewLetVar(
				l.Name(),
				ast.NewApp(
					"partial",
					ast.NewArguments(
						append(
							[]ast.PositionalArgument{ast.NewPositionalArgument(flattened, false)},
							namesToPosArgs(usedNames.slice())...,
						), nil, []interface{}{}),
					f.DebugInfo())))

			names.add(l.Name())
		default:
			panic(fmt.Errorf("Invalid value: %#v", l))
		}
	}

	return append(ss, ast.NewLetFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo()))
}

func flattenInnerStatements(f ast.LetFunction) ast.LetFunction {
	return ast.NewLetFunction(
		f.Name(),
		f.Signature(),
		flattenStatements(f.Lets()),
		f.Body(),
		f.DebugInfo())
}
