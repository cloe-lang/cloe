package desugar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/consts"
	"github.com/coel-lang/coel/src/lib/parse"
)

func TestDesugarEmptyCollection(t *testing.T) {
	for _, s := range []string{
		`(write [])`,
		`(write {})`,
		`(write nil)`,
		`(write [42 []])`,
		`(write (read))`,
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
