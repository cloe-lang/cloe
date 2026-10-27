package match

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
)

func app(f any, args ...any) any {
	return ast.NewPApp(f, args, debug.NewGoInfo(1))
}

func newSwitch(v any, cs []ast.SwitchCase, d any) any {
	if len(cs) == 0 {
		return d
	}

	return ast.NewSwitch(v, cs, d)
}
