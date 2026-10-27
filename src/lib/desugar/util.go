package desugar

import "github.com/cloe-lang/cloe/src/lib/ast"

func signatureToNames(s ast.Signature) names {
	ns := newNames()

	for n := range s.NameToIndex() {
		ns.add(n)
	}

	return ns
}

func prependPositionalsToSig(ns []string, s ast.Signature) ast.Signature {
	return ast.NewSignature(
		append(ns, s.Positionals()...), s.RestPositionals(),
		s.Keywords(), s.RestKeywords())
}

func renamer(a, b string) func(x any) any {
	return func(x any) any {
		return ast.Convert(func(x any) any {
			s, ok := x.(string)

			if !ok || s != a {
				return nil
			}

			return b
		}, x)
	}
}

func reverseSlice(xs []any) []any {
	ys := make([]any, 0, len(xs))

	for i := range xs {
		ys = append(ys, xs[len(xs)-i-1])
	}

	return ys
}
