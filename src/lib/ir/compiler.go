package ir

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/scalar"
)

type compiler struct {
	constantCompiler constantCompiler
	code             []int
	switches         []switchData
	names            []string
	variables        variables
}

func newCompiler(f func(string) (core.Value, error)) compiler {
	return compiler{constantCompiler: newConstantCompiler(f)}
}

func (c *compiler) Compile(ps []string, ls []interface{}, b interface{}) ([]int, []core.Value, []switchData, []string, error) {
	cs, m := c.constantCompiler.Compile(ps, ls, b)

	c.variables = newVariables(m, len(cs))

	for _, p := range ps {
		c.variables.AddNamedVariable(p)
	}

	for _, l := range ls {
		l := l.(ast.LetVar)
		i, err := c.compileExpression(l.Expr())

		if err != nil {
			return nil, nil, nil, nil, err
		}

		c.variables.BindNamedVariable(l.Name(), i)
	}

	_, err := c.compileExpression(b)

	if err != nil {
		return nil, nil, nil, nil, err
	} else if len(c.code) != 0 {
		c.code[len(c.code)-1] = 0
	}

	return c.code, cs, c.switches, c.names, nil
}

func (c *compiler) compileExpression(x interface{}) (int, error) {
	switch x := x.(type) {
	case string:
		return c.variables.GetNamedVariable(x)
	case ast.App:
		a := x.Arguments()

		ps := make([]int, 0, len(a.Positionals()))
		ks := make([]int, 0, len(a.Keywords()))

		for _, p := range a.Positionals() {
			i, err := c.compileExpression(p.Value())

			if err != nil {
				return 0, err
			}

			ps = append(ps, i)
		}

		for _, k := range a.Keywords() {
			i, err := c.compileExpression(k.Value())

			if err != nil {
				return 0, err
			}

			ks = append(ks, i)
		}

		f, err := c.compileExpression(x.Function())

		if err != nil {
			return 0, err
		}

		c.addCode(f)

		c.addCode(len(ps))
		for i, p := range a.Positionals() {
			b := 0

			if p.Expanded() {
				b = 1
			}

			c.addCode(b, ps[i])
		}

		c.addCode(len(ks))
		for i, k := range a.Keywords() {
			if k.Name() == "" {
				c.addCode(expandedKeywordArgument)
			} else {
				c.names = append(c.names, k.Name())
				c.addCode(len(c.names) - 1)
			}

			c.addCode(ks[i])
		}
	case ast.Switch:
		cs := make([]caseData, 0, len(x.Cases()))

		for _, d := range x.Cases() {
			k, err := scalar.Convert(d.Pattern())

			if err != nil {
				return 0, err
			}

			v, err := c.compileExpression(d.Value())

			if err != nil {
				return 0, err
			}

			cs = append(cs, newCaseData(k, v))
		}

		v, err := c.compileExpression(x.Value())

		if err != nil {
			return 0, err
		}

		dc, err := c.compileExpression(x.DefaultCase())

		if err != nil {
			return 0, err
		}

		c.addCode(switchExpression, v, c.addSwitchConstant(newSwitchData(cs, dc)))
	}

	i := c.variables.AddVariable()
	c.addCode(i)
	return i, nil
}

func (c *compiler) addCode(is ...int) {
	c.code = append(c.code, is...)
}

func (c *compiler) addSwitchConstant(s switchData) int {
	c.switches = append(c.switches, s)
	return len(c.switches) - 1
}
