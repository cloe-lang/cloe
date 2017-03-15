package desugar

import "github.com/raviqqe/tisp/src/lib/ast"

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
