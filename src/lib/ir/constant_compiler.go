package ir

import (
	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/core"
)

type constantCompiler struct {
	compileConstant func(s string) (core.Value, error)
	constants       []core.Value
	nameToIndex     map[string]int
}

func newConstantCompiler(f func(string) (core.Value, error)) constantCompiler {
	return constantCompiler{f, nil, map[string]int{}}
}

func (c *constantCompiler) Compile(vs []string, ls []interface{}, b interface{}) ([]core.Value, map[string]int) {
	for _, l := range ls {
		l := l.(ast.LetVar)
		c.compileConstants(l, vs)
		vs = append(vs, l.Name())
	}

	c.compileConstants(b, vs)

	return c.constants, c.nameToIndex
}

func (c *constantCompiler) compileConstants(x interface{}, vs []string) interface{} {
	return ast.Convert(func(x interface{}) interface{} {
		switch x := x.(type) {
		case string:
			for _, s := range vs {
				if x == s {
					return nil
				}
			}

			v, err := c.compileConstant(x)

			if err != nil {
				return nil
			}

			c.addConstant(x, v)
		}

		return nil
	}, x)
}

func (c *constantCompiler) addConstant(s string, v core.Value) {
	c.constants = append(c.constants, v)
	c.nameToIndex[s] = len(c.constants) - 1
}
