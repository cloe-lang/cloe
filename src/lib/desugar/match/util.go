package match

import (
	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
)

func app(f interface{}, args ...interface{}) interface{} {
	return ast.NewPApp(f, args, debug.NewGoInfo(0))
}

func newSwitch(v interface{}, cs []ast.SwitchCase, d interface{}) interface{} {
	if len(cs) == 0 {
		return d
	}

	return ast.NewSwitch(v, cs, d)
}
