package desugar

import "github.com/cloe-lang/cloe/src/lib/ast"

func signatureToNames(s ast.Signature) names {
	ns := newNames()

	for _, n := range s.Names() {
		ns.add(n)
	}

	return ns
}

func prependPositionalsToSig(ns []string, s ast.Signature) ast.Signature {
	return ast.NewSignature(
		append(ns, s.Positionals()...), s.RestPositionals(),
		s.Keywords(), s.RestKeywords())
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
