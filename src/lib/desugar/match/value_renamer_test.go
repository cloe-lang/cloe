package match

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestValueRenamerRename(t *testing.T) {
	r := newValueRenamer(map[string]string{"foo": "bar"})

	for _, x := range []interface{}{
		"foo",
		ast.NewPApp("foo", []interface{}{"bar"}, debug.NewGoInfo(0)),
		ast.NewSwitch("123", []ast.SwitchCase{
			ast.NewSwitchCase("456", "false"),
			ast.NewSwitchCase("123", "true"),
		}, "$matchError"),
		ast.NewApp(
			"foo",
			ast.NewArguments(
				[]ast.PositionalArgument{ast.NewPositionalArgument("123", false)},
				[]ast.KeywordArgument{
					ast.NewKeywordArgument("key", `"value"`),
					ast.NewKeywordArgument("", "bar"),
				}),
			debug.NewGoInfo(0)),
	} {
		r.rename(x)
	}
}

func TestValueRenamerRenameFail(t *testing.T) {
	assert.Panics(t, func() {
		newValueRenamer(map[string]string{"foo": "bar"}).rename(
			ast.NewMatch("nil", []ast.MatchCase{ast.NewMatchCase("nil", "nil")}))
	})
}
