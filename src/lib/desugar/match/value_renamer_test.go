package match

import (
	"testing"

	"github.com/tisp-lang/tisp/src/lib/ast"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

func TestValueRenamerRename(t *testing.T) {
	r := newValueRenamer(map[string]string{"foo": "bar"})

	for _, x := range []interface{}{
		"foo",
		ast.NewPApp("foo", []interface{}{"bar"}, debug.NewGoInfo(0)),
		ast.NewSwitch("123", []ast.SwitchCase{
			ast.NewSwitchCase("456", "false"),
			ast.NewSwitchCase("123", "true"),
		}, `"Error occurred."`),
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
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	newValueRenamer(map[string]string{"foo": "bar"}).rename(
		ast.NewMatch("nil", []ast.MatchCase{ast.NewMatchCase("nil", "nil")}))
}
