package match

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
)

func app(f interface{}, args ...interface{}) interface{} {
	return ast.NewPApp(f, args, debug.NewGoInfo(1))
}

func newSwitch(v interface{}, cs []ast.SwitchCase, d interface{}) interface{} {
	if len(cs) == 0 {
		return d
	}

	return ast.NewSwitch(v, cs, d)
}
