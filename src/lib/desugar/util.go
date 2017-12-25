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

func renamer(a, b string) func(x interface{}) interface{} {
	return func(x interface{}) interface{} {
		return ast.Convert(func(x interface{}) interface{} {
			s, ok := x.(string)

			if !ok || s != a {
				return nil
			}

			return b
		}, x)
	}
}

func reverseSlice(xs []interface{}) []interface{} {
	ys := make([]interface{}, 0, len(xs))

	for i := range xs {
		ys = append(ys, xs[len(xs)-i-1])
	}

	return ys
}
