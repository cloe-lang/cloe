package parse

import (
	"fmt"

	"github.com/raviqqe/tisp/src/lib/ast"
	"github.com/raviqqe/tisp/src/lib/debug"
	"github.com/raviqqe/tisp/src/lib/parse/comb"
)

const (
	commentChar  = ';'
	invalidChars = "\x00"
	quoteString  = "quote"
	spaceChars   = " \t\n\r"
	specialChars = "()[]{}\"'`$"
)

var reserveds = map[string]bool{
	"let":       true,
	"macro":     true,
	quoteString: true,
	"rec":       true,
}

// Parse parses a file into an AST of the language.
func Parse(file, source string) ([]interface{}, error) {
	m, err := newState(file, string(source)).module()()

	if err != nil {
		return nil, err
	}

	return m.([]interface{}), nil
}

func (s *state) module() comb.Parser {
	return s.Exhaust(s.Prefix(s.blank(), s.Many(s.Or(s.let(), s.output()))))
}

func (s *state) let() comb.Parser {
	return s.Lazy(s.strictLet)
}

func (s *state) strictLet() comb.Parser {
	return s.Or(s.letConst(), s.letFunction())
}

func (s *state) letConst() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewLetConst(xs[1].(string), xs[2])
	}, s.list(s.strippedString("let"), s.identifier(), s.expression()))
}

func (s *state) letFunction() comb.Parser {
	return s.withInfo(
		s.list(
			s.strippedString("let"),
			s.list(s.identifier(), s.signature()),
			s.Many(s.let()),
			s.expression()),
		func(x interface{}, i debug.Info) (interface{}, error) {
			xs := x.([]interface{})
			ys := xs[1].([]interface{})
			return ast.NewLetFunction(ys[0].(string), ys[1].(ast.Signature), xs[2].([]interface{}), xs[3], i), nil
		})
}

func (s *state) signature() comb.Parser {
	optArg := s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewOptionalArgument(xs[0].(string), xs[1])
	}, s.strip(s.list(s.identifier(), s.expression())))

	expanded := s.strip(s.expanded(s.identifier()))

	argSet := s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})

		xs0 := xs[0].([]interface{})
		requireds := make([]string, len(xs0))
		for i, r := range xs0 {
			requireds[i] = r.(string)
		}

		xs1 := xs[1].([]interface{})
		optionals := make([]ast.OptionalArgument, len(xs1))
		for i, o := range xs1 {
			optionals[i] = o.(ast.OptionalArgument)
		}

		expanded := ""
		if xs[2] != nil {
			expanded = xs[2].(string)
		}

		return [3]interface{}{requireds, optionals, expanded}
	}, s.And(s.Many(s.identifier()), s.Many(optArg), s.Maybe(expanded)))

	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})

		pas := xs[0].([3]interface{})
		kas, ok := xs[1].([3]interface{})

		if !ok {
			kas = [3]interface{}{[]string{}, []ast.OptionalArgument{}, ""}
		}

		return ast.NewSignature(
			pas[0].([]string), pas[1].([]ast.OptionalArgument), pas[2].(string),
			kas[0].([]string), kas[1].([]ast.OptionalArgument), kas[2].(string),
		)
	}, s.And(argSet, s.Maybe(s.Prefix(s.strippedString("."), argSet))))
}

func (s *state) output() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		expanded := false

		if xs[0] != nil {
			expanded = true
		}

		return ast.NewOutput(xs[1], expanded)
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
		s.firstOrderExpression(),
		s.Lazy(func() comb.Parser { return s.quote(s.expression()) })))
}

func (s *state) firstOrderExpression() comb.Parser {
	return s.Or(
		s.identifier(),
		s.stringLiteral(),
		s.app(),
		s.appFunc("$list", s.sequence("[", "]")),
		s.appFunc("$dict", s.sequence("{", "}")),
		s.appFunc("$set", s.sequence("'{", "}")),
		// s.appFunc("lambda", s.sequence("'(", ")")),
	)
}

func (s *state) quote(p comb.Parser) comb.Parser {
	return s.appQuote(s.Prefix(s.Char('`'), p))
}

func (s *state) appQuote(p comb.Parser) comb.Parser {
	return s.appWithInfo(
		p,
		func(x interface{}) (interface{}, ast.Arguments) {
			return quoteString, ast.NewArguments(
				[]ast.PositionalArgument{ast.NewPositionalArgument(x, false)},
				[]ast.KeywordArgument{},
				[]interface{}{})
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
		func(x interface{}, i debug.Info) (interface{}, error) {
			f, args := f(x)
			return ast.NewApp(f, args, i), nil
		})
}

func (s *state) withInfo(p comb.Parser, f func(interface{}, debug.Info) (interface{}, error)) comb.Parser {
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

		ks := []ast.KeywordArgument{}
		dicts := []interface{}{}

		if xs, ok := xs[1].([2]interface{}); ok {
			ks = xs[0].([]ast.KeywordArgument)
			dicts = xs[1].([]interface{})
		}

		return ast.NewArguments(xs[0].([]ast.PositionalArgument), ks, dicts)
	}, s.And(
		s.positionalArguments(),
		s.Maybe(s.Prefix(s.strippedString("."), s.keywordArguments()))))
}

func (s *state) positionalArguments() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ps := make([]ast.PositionalArgument, len(xs))

		for i, x := range xs {
			ps[i] = x.(ast.PositionalArgument)
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

		ys := xs[0].([]interface{})
		ks := make([]ast.KeywordArgument, len(ys))
		for i, y := range ys {
			ks[i] = y.(ast.KeywordArgument)
		}

		return [2]interface{}{ks, xs[1].([]interface{})}
	}, s.And(s.Many(s.keywordArgument()), s.Many(s.expanded(s.expression()))))
}

func (s *state) keywordArgument() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewKeywordArgument(xs[0].(string), xs[1])
	}, s.And(s.identifier(), s.expression()))
}

func (s *state) identifier() comb.Parser {
	cs := string(commentChar) + invalidChars + spaceChars + specialChars
	p := s.strip(s.Stringify(s.And(s.NotInString(cs+"."), s.Stringify(s.Many(s.NotInString(cs))))))

	return func() (interface{}, error) {
		x, err := p()

		if err != nil {
			return nil, err
		}

		if _, ok := reserveds[x.(string)]; ok {
			return nil, fmt.Errorf("%#v is a reserved identifier.", x)
		}

		return x, nil
	}
}

func (s *state) stringLiteral() comb.Parser {
	c := s.Char('"')

	return s.appQuote(s.Stringify(s.Wrap(
		c,
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		s.strip(c))))
}

func (s *state) list(ps ...comb.Parser) comb.Parser {
	return s.stringWrap("(", s.And(ps...), ")")
}

func (s *state) sequence(l, r string) comb.Parser {
	return s.App(func(x interface{}) interface{} {
		return ast.NewArguments(x.([]ast.PositionalArgument), []ast.KeywordArgument{}, []interface{}{})
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
	return s.Void(s.Many(s.Or(s.InString(spaceChars), s.comment())))
}

func (s *state) comment() comb.Parser {
	return s.Void(s.And(
		s.Char(commentChar),
		s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *state) strippedString(str string) comb.Parser {
	return s.strip(s.String(str))
}
