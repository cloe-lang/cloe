package match

import (
	"fmt"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/gensym"
	"github.com/tisp-lang/tisp/src/lib/scalar"
)

type patternRenamer struct {
	nameMap map[string]string
}

func newPatternRenamer() *patternRenamer {
	return &patternRenamer{map[string]string{}}
}

func (r *patternRenamer) rename(p interface{}) (interface{}, map[string]string) {
	q := r.renameNames(p)
	return q, r.nameMap
}

func (r *patternRenamer) renameNames(p interface{}) interface{} {
	switch x := p.(type) {
	case string:
		if scalar.Defined(x) {
			return x
		}

		r.nameMap[x] = gensym.GenSym(x, "renamed")
		return r.nameMap[x]
	case ast.App:
		switch x.Function().(string) {
		case "$list":
			fallthrough
		case "$dict":
			ps := make([]ast.PositionalArgument, 0, len(x.Arguments().Positionals()))

			for _, p := range x.Arguments().Positionals() {
				ps = append(ps, ast.NewPositionalArgument(r.renameNames(p.Value()), p.Expanded()))
			}

			return ast.NewApp(x.Function(), ast.NewArguments(ps, nil, nil), x.DebugInfo())
		}
	}

	panic(fmt.Errorf("Invalid pattern: %#v", p))
}
