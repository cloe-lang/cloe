package ir

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/scalar"
	"github.com/stretchr/testify/assert"
)

func TestCompilerCompile(t *testing.T) {
	for _, k := range []struct {
		parameters []string
		lets       []interface{}
		body       interface{}
		arguments  []core.Value
		result     core.Value
	}{
		{
			[]string{"x"},
			nil,
			"x",
			[]core.Value{core.Nil},
			core.Nil,
		},
		{
			[]string{},
			nil,
			`"foo"`,
			nil,
			core.NewString("foo"),
		},
		{
			[]string{"f"},
			nil,
			ast.NewPApp("f", []interface{}{"1", "1"}, nil),
			[]core.Value{core.Add},
			core.NewNumber(2),
		},
		{
			[]string{"f", "x", "y"},
			nil,
			ast.NewPApp("f", []interface{}{"x", "y"}, nil),
			[]core.Value{core.Add, core.NewNumber(1), core.NewNumber(2)},
			core.NewNumber(3),
		},
		{
			[]string{"f", "x", "y"},
			nil,
			ast.NewPApp(
				"f",
				[]interface{}{"x", "y", ast.NewPApp("f", []interface{}{"x", "y"}, nil)},
				nil),
			[]core.Value{core.Add, core.NewNumber(1), core.NewNumber(2)},
			core.NewNumber(6),
		},
		{
			[]string{"f", "true"},
			nil,
			ast.NewPApp("f", []interface{}{"true", "1"}, nil),
			[]core.Value{core.Add, core.NewNumber(1)},
			core.NewNumber(2),
		},
		{
			[]string{"f", "true"},
			nil,
			ast.NewPApp("f", []interface{}{"true"}, nil),
			[]core.Value{core.Add, core.NewNumber(1)},
			core.NewNumber(1),
		},
		{
			[]string{"f", "x"},
			[]interface{}{ast.NewLetVar("true", "x")},
			ast.NewPApp("f", []interface{}{"true"}, nil),
			[]core.Value{core.Add, core.NewNumber(1)},
			core.NewNumber(1),
		},
		{
			[]string{"f", "x"},
			nil,
			ast.NewApp(
				"f",
				ast.NewArguments([]ast.PositionalArgument{ast.NewPositionalArgument("x", true)}, nil),
				nil),
			[]core.Value{core.Add, core.NewList(core.NewNumber(1), core.NewNumber(2))},
			core.NewNumber(3),
		},
		{
			[]string{"f", "x"},
			nil,
			ast.NewApp(
				"f",
				ast.NewArguments(nil, []ast.KeywordArgument{ast.NewKeywordArgument("foo", "x")}),
				nil),
			[]core.Value{testFunction, core.NewNumber(42)},
			core.NewNumber(42),
		},
		{
			[]string{"f", "x"},
			nil,
			ast.NewApp(
				"f",
				ast.NewArguments(nil, []ast.KeywordArgument{ast.NewKeywordArgument("", "x")}),
				nil),
			[]core.Value{
				testFunction,
				core.NewDictionary([]core.KeyValue{{Key: core.NewString("foo"), Value: core.NewNumber(42)}}),
			},
			core.NewNumber(42),
		},
		{
			[]string{"x", "y"},
			nil,
			ast.NewSwitch("x", []ast.SwitchCase{ast.NewSwitchCase(`"foo"`, "y")}, "x"),
			[]core.Value{core.NewString("foo"), core.NewNumber(42)},
			core.NewNumber(42),
		},
		{
			[]string{"add", "rest", "list"},
			nil,
			ast.NewApp(
				"add",
				ast.NewArguments([]ast.PositionalArgument{
					ast.NewPositionalArgument(ast.NewPApp("rest", []interface{}{ast.NewPApp("rest", []interface{}{"list"}, nil)}, nil), true),
				}, nil),
				nil),
			[]core.Value{
				core.Add,
				core.Rest,
				core.NewList(core.NewNumber(1), core.NewNumber(1), core.NewNumber(1)),
			},
			core.NewNumber(1),
		},
	} {
		c := newCompiler(scalar.Convert)
		bs, cs, scs, ns, err := c.Compile(k.parameters, k.lets, k.body)
		assert.Nil(t, err)

		i := NewInterpreter(bs, scs, ns, append(cs, k.arguments...))
		b, e := core.EvalBoolean(core.PApp(core.Equal, k.result, i.Interpret()))

		assert.Nil(t, e)
		assert.True(t, bool(b))
	}
}

func TestCompilerCompileError(t *testing.T) {
	for _, k := range []struct {
		parameters []string
		lets       []interface{}
		body       interface{}
	}{
		{
			[]string{"x"},
			nil,
			"y",
		},
	} {
		c := newCompiler(scalar.Convert)
		_, _, _, _, err := c.Compile(k.parameters, k.lets, k.body)

		assert.NotNil(t, err)
	}
}
