package desugar

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func TestDesugarAnonymousFunctions(t *testing.T) {
	for _, s := range []interface{}{
		ast.NewLetVar(
			"foo",
			ast.NewAnonymousFunction(ast.NewSignature(nil, nil, "", nil, nil, ""), "123")),
		ast.NewLetFunction(
			"foo",
			ast.NewSignature(nil, nil, "", nil, nil, ""),
			nil,
			ast.NewAnonymousFunction(ast.NewSignature(nil, nil, "", nil, nil, ""), "123"),
			debug.NewGoInfo(0)),
		ast.NewLetFunction(
			"foo",
			ast.NewSignature(nil, nil, "", nil, nil, ""),
			[]interface{}{
				ast.NewLetVar(
					"x",
					ast.NewAnonymousFunction(ast.NewSignature(nil, nil, "", nil, nil, ""), "123")),
			},
			"x",
			debug.NewGoInfo(0)),
	} {
		t.Logf("%#v", s)

		for _, s := range desugarAnonymousFunctions(s) {
			ast.Convert(func(x interface{}) interface{} {
				if _, ok := x.(ast.AnonymousFunction); ok {
					t.Fail()
				}

				return nil
			}, s)
		}
	}
}
