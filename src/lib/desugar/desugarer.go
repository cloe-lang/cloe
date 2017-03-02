package desugar

import (
	"github.com/raviqqe/tisp/src/lib/ast"
	"log"
)

type desugarer struct{}

func newDesugarer() desugarer {
	return desugarer{}
}

func (d *desugarer) desugar(module []interface{}) []interface{} {
	ss := make([]interface{}, 0, 2*len(module)) // TODO: Best cap?

	for _, s := range module {
		ss = append(ss, d.desugarStatement(s)...)
	}

	return ss
}

func (d *desugarer) desugarStatement(s interface{}) []interface{} {
	switch s := s.(type) {
	case ast.LetFunction:
		return d.desugarLetFunction(s)
	default:
		return []interface{}{s}
	}
}

func (d *desugarer) desugarLetFunction(f ast.LetFunction) []interface{} {
	cap := len(f.Lets()) + 1
	ss := make([]interface{}, 0, cap)
	ls := make([]interface{}, 0, cap)

	for _, l := range f.Lets() {
		switch l := l.(type) {
		case ast.LetConst:
			ls = append(ls, l)
		case ast.LetFunction:
			unnested := "$" + f.Name() + "$" + l.Name()
			// TODO: Remove names in signatureToNames(l.Signature()).
			usedNames := signatureToNames(f.Signature()).find(l.Body()).slice()

			ss = append(ss, ast.NewLetFunction(
				unnested,
				prependPosReqsToSig(l.Signature(), usedNames),
				l.Lets(),
				l.Body()))

			ls = append(ls, ast.NewLetConst(l.Name(), ast.NewApp("partial", ast.NewArguments(
				append(
					[]ast.PositionalArgument{ast.NewPositionalArgument(unnested, false)},
					namesToPosArgs(usedNames)...,
				), []ast.KeywordArgument{}, []interface{}{}))))
		default:
			log.Panicf("Invalid value: %#v\n", l)
		}
	}

	return append(ss, ast.NewLetFunction(f.Name(), f.Signature(), ls, f.Body()))
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
