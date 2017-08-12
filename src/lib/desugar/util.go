package desugar

import "github.com/tisp-lang/tisp/src/lib/ast"

func signatureToNames(s ast.Signature) names {
	ns := newNames(append(s.PosReqs(), s.KeyReqs()...)...)

	for _, o := range s.PosOpts() {
		ns.add(o.Name())
	}

	if r := s.PosRest(); r != "" {
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

func prependPosReqsToSig(ns []string, s ast.Signature) ast.Signature {
	return ast.NewSignature(
		append(ns, s.PosReqs()...), s.PosOpts(), s.PosRest(),
		s.KeyReqs(), s.KeyOpts(), s.KeyRest())
}
