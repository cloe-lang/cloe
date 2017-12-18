package match

import (
	"testing"

	"github.com/coel-lang/coel/src/lib/ast"
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/stretchr/testify/assert"
)

func TestValueRenamerRename(t *testing.T) {
	r := newValueRenamer(map[string]string{"foo": "bar"})

	for _, x := range []interface{}{
		"foo",
		ast.NewPApp("foo", []interface{}{"bar"}, debug.NewGoInfo(0)),
		ast.NewMatch("123", []ast.MatchCase{
			ast.NewMatchCase("456", "false"),
			ast.NewMatchCase("123", "true"),
		}),
		ast.NewApp(
			"foo",
			ast.NewArguments(
				[]ast.PositionalArgument{ast.NewPositionalArgument("123", false)},
				[]ast.KeywordArgument{ast.NewKeywordArgument("key", `"value"`)},
				[]interface{}{"bar"}),
			debug.NewGoInfo(0)),
	} {
		r.rename(x)
	}
}

func TestValueRenamerRenameFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	newValueRenamer(map[string]string{"foo": "bar"}).rename(
		ast.NewSwitch("nil", []ast.SwitchCase{ast.NewSwitchCase("nil", "nil")}, `"not so match"`))
}
