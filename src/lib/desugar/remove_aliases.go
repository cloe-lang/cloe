package desugar

import (
	"github.com/coel-lang/coel/src/lib/ast"
)

func removeAliases(x interface{}) []interface{} {
	return []interface{}{ast.Convert(func(x interface{}) interface{} {
		f, ok := x.(ast.DefFunction)

		if !ok {
			return nil
		}

		ls := make([]interface{}, 0, len(f.Lets()))
		b := f.Body()

		for _, l := range reverseSlice(f.Lets()) {
			l := l.(ast.LetVar)
			s, ok := l.Expr().(string)

			if !ok {
				ls = append(ls, l)
				continue
			}

			r := renamer(l.Name(), s)

			for i, l := range ls {
				l := l.(ast.LetVar)
				ls[i] = ast.NewLetVar(l.Name(), r(l.Expr()))
			}

			b = r(b)
		}

		return ast.NewDefFunction(f.Name(), f.Signature(), reverseSlice(ls), b, f.DebugInfo())
	}, x)}
}
