package parse

import (
	"fmt"
	"strconv"

	"github.com/cloe-lang/cloe/src/lib/ast"
	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/cloe-lang/cloe/src/lib/debug"
	"github.com/cloe-lang/cloe/src/lib/parse/comb"
)

const (
	defString       = "def"
	commentChar     = ';'
	importString    = "import"
	invalidChars    = "\x00"
	letString       = "let"
	mutualRecString = "mr"
	matchString     = "match"
	spaceChars      = " \t\n\r"
	specialChars    = "()[]{}\"\\$"
)

var reserveds = map[string]bool{
	defString:       true,
	importString:    true,
	letString:       true,
	matchString:     true,
	mutualRecString: true,
}

// MainModule parses a main module file into an AST.
func MainModule(file, source string) ([]interface{}, error) {
	m, err := newState(file, source).mainModule()()

	if err != nil {
		return nil, err
	}

	return m.([]interface{}), nil
}

// SubModule parses a sub module file into an AST.
func SubModule(file, source string) ([]interface{}, error) {
	m, err := newState(file, source).subModule()()

	if err != nil {
		return nil, err
	}

	return m.([]interface{}), nil
}

func (s *state) mainModule() comb.Parser {
	return s.Prefix(s.Maybe(s.shebang()), s.module(s.importModule(), s.let(), s.effect()))
}

func (s *state) subModule() comb.Parser {
	return s.module(s.importModule(), s.let())
}

func (s *state) module(ps ...comb.Parser) comb.Parser {
	return s.exhaust(s.Prefix(s.blank(), s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return append(xs[0].([]interface{}), xs[1].([]interface{})...)
	}, s.And(s.Many(s.importModule()), s.Many(s.Or(ps...))))))
}

func (s *state) importModule() comb.Parser {
	return s.withInfo(
		s.list(s.strippedString(importString), s.stringLiteral()),
		func(x interface{}, i *debug.Info) (interface{}, error) {
			s, err := strconv.Unquote(x.([]interface{})[1].(string))

			if err != nil {
				return nil, err
			}

			return ast.NewImport(s, i), nil
		})
}

func (s *state) let() comb.Parser {
	return s.Lazy(s.strictLet)
}

func (s *state) strictLet() comb.Parser {
	return s.Or(s.letVar(), s.letMatch(), s.letFunction(), s.mutuallyRecursiveDefFunctions())
}

func (s *state) letVar() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewLetVar(xs[1].(string), xs[2])
	}, s.list(s.strippedString(letString), s.identifier(), s.expression()))
}

func (s *state) letMatch() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewLetMatch(xs[1], xs[2])
	}, s.list(s.strippedString(letString), s.Or(s.listLiteral(), s.dictLiteral()), s.expression()))
}

func (s *state) letFunction() comb.Parser {
	return s.withInfo(
		s.list(
			s.strippedString(defString),
			s.list(s.identifier(), s.signature()),
			s.Many(s.let()),
			s.expression()),
		func(x interface{}, i *debug.Info) (interface{}, error) {
			xs := x.([]interface{})
			ys := xs[1].([]interface{})
			return ast.NewDefFunction(ys[0].(string), ys[1].(ast.Signature), xs[2].([]interface{}), xs[3], i), nil
		})
}

func (s *state) optionalParameter() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewOptionalParameter(xs[0].(string), xs[1])
	}, s.strip(s.And(s.identifier(), s.expression())))
}

func (s *state) expandedArgument() comb.Parser {
	return s.strip(s.expanded(s.identifier()))
}

func (s *state) positionalParameters() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ys := xs[0].([]interface{})

		ss := make([]string, 0, len(ys))

		for _, y := range ys {
			ss = append(ss, y.(string))
		}

		s := ""

		if xs[1] != nil {
			s = xs[1].(string)
		}

		return [2]interface{}{ss, s}
	}, s.And(s.Many(s.identifier()), s.Maybe(s.expandedArgument())))
}

func (s *state) keywordParameters() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ys := xs[0].([]interface{})

		os := make([]ast.OptionalParameter, 0, len(ys))

		for _, y := range ys {
			os = append(os, y.(ast.OptionalParameter))
		}

		s := ""

		if xs[1] != nil {
			s = xs[1].(string)
		}

		return [2]interface{}{os, s}
	}, s.And(s.Many(s.optionalParameter()), s.Maybe(s.expandedArgument())))
}

func (s *state) signature() comb.Parser {

	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})

		pas := xs[0].([2]interface{})
		kas, ok := xs[1].([2]interface{})

		if !ok {
			kas = [2]interface{}{[]ast.OptionalParameter(nil), ""}
		}

		return ast.NewSignature(
			pas[0].([]string), pas[1].(string),
			kas[0].([]ast.OptionalParameter), kas[1].(string),
		)
	}, s.And(s.positionalParameters(), s.Maybe(s.Prefix(s.strippedString("."), s.keywordParameters()))))
}

func (s *state) effect() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		expanded := false

		if xs[0] != nil {
			expanded = true
		}

		return ast.NewEffect(xs[1], expanded)
	}, s.And(s.Maybe(s.String("..")), s.expression()))
}

func (s *state) expanded(p comb.Parser) comb.Parser {
	return s.Prefix(s.String(".."), p)
}

func (s *state) expression() comb.Parser {
	return s.Lazy(s.strictExpression)
}

func (s *state) strictExpression() comb.Parser {
	return s.strip(s.Or(
		s.identifier(),
		s.stringLiteral(),
		s.match(),
		s.app(),
		s.listLiteral(),
		s.dictLiteral(),
		s.anonymousFunction(),
	))
}

func (s *state) listLiteral() comb.Parser {
	return s.appFunc(consts.Names.ListFunction, s.sequence("[", "]"))
}

func (s *state) dictLiteral() comb.Parser {
	return s.appFunc(consts.Names.DictionaryFunction, s.sequence("{", "}"))
}

func (s *state) anonymousFunction() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewAnonymousFunction(xs[1].([]interface{})[0].(ast.Signature), xs[2])
	}, s.list(s.strippedString("\\"), s.list(s.signature()), s.expression()))
}

func (s *state) match() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ks := xs[2].([]interface{})

		cs := make([]ast.MatchCase, 0, len(ks))

		for _, k := range ks {
			xs := k.([]interface{})
			cs = append(cs, ast.NewMatchCase(xs[0], xs[1]))
		}

		return ast.NewMatch(xs[1], cs)
	}, s.list(
		s.strippedString(matchString),
		s.expression(),
		s.Many1(s.And(s.pattern(), s.expression()))))
}

func (s *state) pattern() comb.Parser {
	return s.Or(
		s.identifier(),
		s.stringLiteral(),
		s.listLiteral(),
		s.dictLiteral())
}

func (s *state) mutuallyRecursiveDefFunctions() comb.Parser {
	return s.withInfo(
		s.list(s.strippedString(mutualRecString), s.Many(s.letFunction())),
		func(x interface{}, i *debug.Info) (interface{}, error) {
			xs := x.([]interface{})[1].([]interface{})
			fs := make([]ast.DefFunction, 0, len(xs))

			for _, l := range xs {
				fs = append(fs, l.(ast.DefFunction))
			}

			return ast.NewMutualRecursion(fs, i), nil
		})
}

func (s *state) app() comb.Parser {
	return s.appWithInfo(
		s.list(s.expression(), s.arguments()),
		func(x interface{}) (interface{}, ast.Arguments) {
			xs := x.([]interface{})
			return xs[0], xs[1].(ast.Arguments)
		})
}

func (s *state) appWithInfo(p comb.Parser, f func(interface{}) (interface{}, ast.Arguments)) comb.Parser {
	return s.withInfo(
		p,
		func(x interface{}, i *debug.Info) (interface{}, error) {
			f, args := f(x)
			return ast.NewApp(f, args, i), nil
		})
}

func (s *state) withInfo(p comb.Parser, f func(interface{}, *debug.Info) (interface{}, error)) comb.Parser {
	return func() (interface{}, error) {
		i := s.debugInfo()
		x, err := p()

		if err != nil {
			return nil, err
		}

		return f(x, i)
	}
}

func (s *state) arguments() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})

		ks := []ast.KeywordArgument(nil)

		if xs[1] != nil {
			ks = xs[1].([]ast.KeywordArgument)
		}

		return ast.NewArguments(xs[0].([]ast.PositionalArgument), ks)
	}, s.And(
		s.positionalArguments(),
		s.Maybe(s.Prefix(s.strippedString("."), s.keywordArguments()))))
}

func (s *state) positionalArguments() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ps := make([]ast.PositionalArgument, 0, len(xs))

		for _, x := range xs {
			ps = append(ps, x.(ast.PositionalArgument))
		}

		return ps
	}, s.Many(s.positionalArgument()))
}

func (s *state) positionalArgument() comb.Parser {
	unexpanded := s.App(func(x interface{}) interface{} {
		return ast.NewPositionalArgument(x, false)
	}, s.expression())

	expanded := s.App(func(x interface{}) interface{} {
		return ast.NewPositionalArgument(x, true)
	}, s.expanded(s.expression()))

	return s.Or(unexpanded, expanded)
}

func (s *state) keywordArguments() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ks := make([]ast.KeywordArgument, 0, len(xs))

		for _, x := range xs {
			switch x := x.(type) {
			case ast.KeywordArgument:
				ks = append(ks, x)
			default:
				ks = append(ks, ast.NewKeywordArgument("", x))
			}
		}

		return ks
	}, s.Many(s.Or(s.keywordArgument(), s.expanded(s.expression()))))
}

func (s *state) keywordArgument() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewKeywordArgument(xs[0].(string), xs[1])
	}, s.And(s.identifier(), s.expression()))
}

func (s *state) identifier() comb.Parser {
	cs := string(commentChar) + invalidChars + spaceChars + specialChars
	p := s.strip(s.Stringify(s.And(s.NotChars(cs+"."), s.Stringify(s.Many(s.NotChars(cs))))))

	return func() (interface{}, error) {
		x, err := p()

		if err != nil {
			return nil, err
		}

		if _, ok := reserveds[x.(string)]; ok {
			return nil, fmt.Errorf("%#v is a reserved identifier", x)
		}

		return x, nil
	}
}

func (s *state) stringLiteral() comb.Parser {
	c := s.Char('"')

	return s.Stringify(s.And(
		c,
		s.Many(s.Or(
			s.NotChars("\"\\"),
			s.String("\\\""),
			s.String("\\\\"),
			s.String("\\n"),
			s.String("\\t"),
		)),
		s.strip(c)))
}

func (s *state) list(ps ...comb.Parser) comb.Parser {
	return s.stringWrap("(", s.And(ps...), ")")
}

func (s *state) sequence(l, r string) comb.Parser {
	return s.App(func(x interface{}) interface{} {
		return ast.NewArguments(x.([]ast.PositionalArgument), nil)
	}, s.stringWrap(l, s.positionalArguments(), r))
}

func (s *state) stringWrap(l string, p comb.Parser, r string) comb.Parser {
	return s.Wrap(s.strippedString(l), p, s.strippedString(r))
}

func (s *state) appFunc(ident string, p comb.Parser) comb.Parser {
	return s.appWithInfo(
		p,
		func(x interface{}) (interface{}, ast.Arguments) {
			return ident, x.(ast.Arguments)
		})
}

func (s *state) strip(p comb.Parser) comb.Parser {
	b := s.blank()
	return s.Wrap(s.None(), p, b)
}

func (s *state) blank() comb.Parser {
	return s.Void(s.Many(s.Or(s.Chars(spaceChars), s.comment())))
}

func (s *state) comment() comb.Parser {
	return s.Void(s.And(
		s.Char(commentChar),
		s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *state) strippedString(str string) comb.Parser {
	return s.strip(s.String(str))
}

func (s *state) exhaust(p comb.Parser) comb.Parser {
	return s.Exhaust(p, func(s comb.State) error {
		return fmt.Errorf("SyntaxError:%d:%d:\t%s", s.LineNumber(), s.LinePosition(), s.Line())
	})
}

func (s *state) shebang() comb.Parser {
	return s.Wrap(s.String("#!"), s.Many(s.NotChar('\n')), s.String("\n"))
}
