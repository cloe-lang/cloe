package parse

import (
	"../ast"
	"./comb"
)

const (
	commentChar  = ';'
	invalidChars = "\x00"
	quoteString  = "quote"
	spaceChars   = " \t\n\r"
	specialChars = "()[]{}\"'`$"
)

func Parse(source string) []interface{} {
	m, err := newState(source).module()()

	if err != nil {
		panic(err.Error())
	}

	return m.([]interface{})
}

func (s *state) module() comb.Parser {
	return s.Exhaust(s.Wrap(s.blank(), s.Many(s.Or(s.letConst(), s.letFunction(), s.output())), s.None()))
}

func (s *state) letConst() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		return ast.NewLetConst(xs[1].(string), xs[2])
	}, s.list(s.strippedString("let"), s.identifier(), s.expression()))
}

func (s *state) letFunction() comb.Parser {
	return s.App(func(x interface{}) interface{} {
		xs := x.([]interface{})
		ys := xs[1].([]interface{})
		return ast.NewLetFunction(ys[0].(string), ys[1].(ast.Signature), xs[2])
	}, s.list(s.strippedString("let"), s.list(s.identifier(), s.signature()), s.expression()))
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
	}, s.And(argSet, s.Maybe(s.Wrap(s.strippedString("."), argSet, s.None()))))
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
	return s.Wrap(s.String(".."), p, s.None())
}

func (s *state) expressions() comb.Parser {
	return s.Lazy(s.strictExpressions)
}

func (s *state) strictExpressions() comb.Parser {
	return s.Many(s.expression())
}

func (s *state) expression() comb.Parser {
	return s.strip(s.Or(
		s.firstOrderExpression(),
		s.Lazy(func() comb.Parser { return s.quote(s.expression()) })))
}

func (s *state) firstOrderExpression() comb.Parser {
	return s.Or(
		s.identifier(),
		s.String(".."),
		s.String("."),
		s.stringLiteral(),
		s.sequence("(", ")"),
		s.prepend("list", s.sequence("[", "]")),
		s.prepend("dict", s.sequence("{", "}")),
		s.prepend("set", s.sequence("'{", "}")),
		s.prepend("lambda", s.sequence("'(", ")")))
}

func (s *state) quote(p comb.Parser) comb.Parser {
	return s.And(s.Replace(quoteString, s.Char('`')), p)
}

func (s *state) identifier() comb.Parser {
	cs := string(commentChar) + invalidChars + spaceChars + specialChars
	return s.strip(s.Stringify(s.And(s.NotInString(cs+"."), s.Stringify(s.Many(s.NotInString(cs))))))
}

func (s *state) stringLiteral() comb.Parser {
	c := s.Char('"')
	f := func(x interface{}) interface{} {
		return []interface{}{quoteString, x}
	}

	return s.App(f, s.Stringify(s.Wrap(
		c,
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		s.strip(c))))
}

func (s *state) list(ps ...comb.Parser) comb.Parser {
	return s.Wrap(s.strippedString("("), s.And(ps...), s.strippedString(")"))
}

func (s *state) sequence(l, r string) comb.Parser {
	return s.Wrap(s.strippedString(l), s.expressions(), s.strippedString(r))
}

func (s *state) prepend(x interface{}, p comb.Parser) comb.Parser {
	return s.App(func(any interface{}) interface{} {
		return append([]interface{}{x}, any.([]interface{})...)
	}, p)
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
