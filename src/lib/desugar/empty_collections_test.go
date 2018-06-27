package desugar

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/parse"
	"github.com/stretchr/testify/assert"
)

func TestDesugarEmptyCollection(t *testing.T) {
	for _, s := range []string{
		`(print [])`,
		`(print {})`,
		`(print nil)`,
		`(print [42 []])`,
		`(print (read))`,
	} {
		m, err := parse.MainModule("<test>", s)
		assert.Nil(t, err)

		for _, s := range m {
			for _, s := range desugarEmptyCollection(s) {
				ast.Convert(func(x interface{}) interface{} {
					a, ok := x.(ast.App)

					if !ok {
						return nil
					}

					switch a.Function() {
					case consts.Names.EmptyList, consts.Names.EmptyDictionary:
						assert.NotZero(t, len(a.Arguments().Positionals()))
					}

					return nil
				}, s)
			}
		}
	}
}
