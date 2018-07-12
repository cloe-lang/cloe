package compile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/desugar"
	"github.com/cloe-lang/cloe/src/lib/ir"
	"github.com/cloe-lang/cloe/src/lib/modules"
	"github.com/cloe-lang/cloe/src/lib/parse"
)

type compiler struct {
	env   environment
	cache modulesCache
}

func newCompiler(e environment, c modulesCache) compiler {
	return compiler{e, c}
}

func (c *compiler) compileModule(m []interface{}, d string) ([]Effect, error) {
	es := []Effect{}

	for _, s := range m {
		switch x := s.(type) {
		case ast.LetVar:
			v, err := c.expressionToValue(x.Expr())

			if err != nil {
				return nil, err
			}

			c.env.set(x.Name(), v)
		case ast.DefFunction:
			s, err := c.compileSignature(x.Signature())

			if err != nil {
				return nil, err
			}

			f, err := ir.CreateFunction(s, x.Signature().Names(), x.Lets(), x.Body(), c.env.get)

			if err != nil {
				return nil, err
			}

			c.env.set(x.Name(), f)
		case ast.Effect:
			v, err := c.expressionToValue(x.Expr())

			if err != nil {
				return nil, err
			}

			es = append(es, NewEffect(v, x.Expanded()))
		case ast.Import:
			if c.cache == nil {
				return nil, errors.New("import statement is unavailable")
			}

			m, ok := modules.Modules[x.Path()]

			if !ok {
				var err error
				m, err = c.importLocalModule(x.Path(), d)

				if err != nil {
					return nil, err
				}
			}

			p := x.Prefix()

			if p != "" {
				p += "."
			}

			for k, v := range m {
				c.env.set(p+k, v)
			}
		default:
			panic(fmt.Errorf("Invalid type: %#v", x))
		}
	}

	return es, nil
}

func (c *compiler) expressionToValue(x interface{}) (core.Value, error) {
	switch x := x.(type) {
	case string:
		return c.env.get(x)
	case ast.App:
		a := x.Arguments()
		ps := make([]core.PositionalArgument, 0, len(a.Positionals()))
		ks := make([]core.KeywordArgument, 0, len(a.Keywords()))

		for _, p := range a.Positionals() {
			v, err := c.expressionToValue(p.Value())

			if err != nil {
				return nil, err
			}

			ps = append(ps, core.NewPositionalArgument(v, p.Expanded()))
		}

		for _, k := range a.Keywords() {
			v, err := c.expressionToValue(k.Value())

			if err != nil {
				return nil, err
			}

			ks = append(ks, core.NewKeywordArgument(k.Name(), v))
		}

		f, err := c.expressionToValue(x.Function())

		if err != nil {
			return nil, err
		}

		return core.AppWithInfo(f, core.NewArguments(ps, ks), x.DebugInfo()), nil
	case ast.Switch:
		kvs := make([]core.KeyValue, 0, len(x.Cases()))

		for _, cs := range x.Cases() {
			k, err := c.env.get(cs.Pattern())

			if err != nil {
				return nil, err
			}

			v, err := c.expressionToValue(cs.Value())

			if err != nil {
				return nil, err
			}

			kvs = append(kvs, core.KeyValue{Key: k, Value: v})
		}

		v, err := c.expressionToValue(x.Value())

		if err != nil {
			return nil, err
		}

		d, err := c.expressionToValue(core.NewDictionary(kvs))

		if err != nil {
			return nil, err
		}

		dc, err := c.expressionToValue(x.DefaultCase())

		if err != nil {
			return nil, err
		}

		return core.PApp(
			core.NewLazyFunction(
				core.NewSignature(nil, "", nil, ""),
				func(...core.Value) core.Value {
					b, err := core.EvalBoolean(core.PApp(core.Include, d, v))

					if err != nil {
						return err
					} else if !b {
						return dc
					}

					return core.PApp(core.Index, d, v)
				})), nil
	}

	panic(fmt.Sprintf("Invalid type: %#v", x))
}

func (c *compiler) compileSignature(s ast.Signature) (core.Signature, error) {
	os := make([]core.OptionalParameter, 0, len(s.Keywords()))

	for _, o := range s.Keywords() {
		v, err := c.expressionToValue(o.DefaultValue())

		if err != nil {
			return core.Signature{}, err
		}

		os = append(os, core.NewOptionalParameter(o.Name(), v))
	}

	return core.NewSignature(s.Positionals(), s.RestPositionals(), os, s.RestKeywords()), nil
}

func (c *compiler) importLocalModule(p, d string) (module, error) {
	var err error

	if p[0] == '.' && (p[:2] == "./" || p[:3] == "../") {
		p, err = filepath.Abs(filepath.FromSlash(path.Join(d, p)))

		if err != nil {
			return nil, err
		}
	} else if !path.IsAbs(p) {
		d, err := consts.GetModulesDirectory()

		if err != nil {
			return nil, err
		}

		p = path.Join(d, p)
	}

	if m, ok := c.cache.Get(p); ok {
		return m, nil
	}

	m, err := c.compileSubModule(p)

	if err != nil {
		return nil, err
	}

	if err := c.cache.Set(p, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *compiler) compileSubModule(p string) (module, error) {
	if i, err := os.Stat(p); err == nil && i.IsDir() {
		p = path.Join(p, consts.ModuleFilename)
	}

	bs, err := ioutil.ReadFile(filepath.FromSlash(p + consts.FileExtension))

	if err != nil {
		return nil, err
	}

	m, err := parse.SubModule(p, string(bs))

	if err != nil {
		return nil, err
	}

	cc := newCompiler(builtinsEnvironment(), c.cache)
	c = &cc
	_, err = c.compileModule(desugar.Desugar(m), path.Dir(p))

	if err != nil {
		return nil, err
	}

	return c.env.toMap(), nil
}
