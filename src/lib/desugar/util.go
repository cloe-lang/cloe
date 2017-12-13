package desugar

import "github.com/coel-lang/coel/src/lib/ast"

func signatureToNames(s ast.Signature) names {
	ns := newNames()

	for n := range s.NameToIndex() {
		ns.add(n)
	}

	return ns
}

func prependPosReqsToSig(ns []string, s ast.Signature) ast.Signature {
	return ast.NewSignature(
		append(ns, s.PosReqs()...), s.PosOpts(), s.PosRest(),
		s.KeyReqs(), s.KeyOpts(), s.KeyRest())
}
