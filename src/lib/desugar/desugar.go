package desugar

import (
	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/util"
)

// Desugar desugars a module of statements in AST.
func Desugar(module []interface{}) []interface{} {
	return desugarStatements(module)
}

func desugarStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.LetFunction:
		return desugarLetFunction(s)
	default:
		return []interface{}{s}
	}
}

func desugarStatements(old []interface{}) []interface{} {
	new := make([]interface{}, 0, 2*len(old)) // TODO: Best cap?

	for _, s := range old {
		new = append(new, desugarStatement(s)...)
	}

	return new
}

func desugarLetFunction(f ast.LetFunction) []interface{} {
	f = desugarInnerStatements(f)

	ss := make([]interface{}, 0)
	ls := make([]interface{}, 0)
	names := signatureToNames(f.Signature())

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetConst:
			ls = append(ls, l)
			names.add(l.Name())
		case ast.LetFunction:
			unnested := f.Name() + "$" + l.Name()

			usedNames := names.find(l.Body())
			for _, l := range l.Lets() {
				usedNames.merge(names.find(l.(ast.LetConst)))
			}
			usedNames.subtract(signatureToNames(l.Signature()))

			ss = append(ss, ast.NewLetFunction(
				unnested,
				prependPosReqsToSig(l.Signature(), usedNames.slice()),
				l.Lets(),
				l.Body(),
				l.DebugInfo()))

			ls = append(ls, ast.NewLetConst(
				l.Name(),
				ast.NewApp(
					"partial",
					ast.NewArguments(
						append(
							[]ast.PositionalArgument{ast.NewPositionalArgument(unnested, false)},
							namesToPosArgs(usedNames.slice())...,
						), []ast.KeywordArgument{}, []interface{}{}),
					f.DebugInfo())))

			names.add(l.Name())
		default:
			util.Fail("Invalid value: %#v\n", l)
		}
	}

	return append(ss, ast.NewLetFunction(f.Name(), f.Signature(), ls, f.Body(), f.DebugInfo()))
}

func desugarInnerStatements(f ast.LetFunction) ast.LetFunction {
	return ast.NewLetFunction(
		f.Name(),
		f.Signature(),
		desugarStatements(f.Lets()),
		f.Body(),
		f.DebugInfo())
}

func signatureToNames(s ast.Signature) names {
	ns := newNames()

	for _, r := range s.PosReqs() {
		ns.add(r)
	}

	for _, o := range s.PosOpts() {
		ns.add(o.Name())
	}

	if r := s.PosRest(); r != "" {
		ns.add(r)
	}

	for _, r := range s.KeyReqs() {
		ns.add(r)
	}

	for _, o := range s.KeyOpts() {
		ns.add(o.Name())
	}

	if r := s.KeyRest(); r != "" {
		ns.add(r)
	}

	return ns
}

func namesToPosArgs(ns []string) []ast.PositionalArgument {
	ps := make([]ast.PositionalArgument, len(ns))

	for i, n := range ns {
		ps[i] = ast.NewPositionalArgument(n, false)
	}

	return ps
}

func prependPosReqsToSig(s ast.Signature, ns []string) ast.Signature {
	return ast.NewSignature(
		append(ns, s.PosReqs()...), s.PosOpts(), s.PosRest(),
		s.KeyReqs(), s.KeyOpts(), s.KeyRest())
}
