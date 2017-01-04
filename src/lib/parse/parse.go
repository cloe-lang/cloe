package parse

import (
	"../types"
	"./comb"
)

const SPACE_CHARS = " ,\t\n\r"

func Parse(source string) types.Object {
	o, err := newState(source).module()()

	if err != nil {
		panic(err.Error())
	}

	return o
}

func (s *state) module() comb.Parser {
	return s.Exhaust(s.elems())
}

func (s *state) elems() comb.Parser {
	return s.Lazy(s.strictElems)
}

func (s *state) strictElems() comb.Parser {
	return s.Many(s.elem())
}

func (s *state) elem() comb.Parser {
	ps := []comb.Parser{s.atom(), s.list(), s.array(), s.dict()}

	return s.strip(s.Or(append(ps, s.quotes(ps...)...)...))
}

func (s *state) atom() comb.Parser {
	return s.Or(s.stringLiteral(), s.identifier())
}

func (s *state) identifier() comb.Parser {
	return s.stringify(s.Many1(s.NotInString("()[]{}$'\x00" + SPACE_CHARS)))
}

func (s *state) stringLiteral() comb.Parser {
	b := s.blank()
	c := s.Char('"')

	return s.stringify(s.Wrap(
		c,
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		s.And(c, b)))
}

func (s *state) list() comb.Parser {
	return s.sequence('(', ')')
}

func (s *state) array() comb.Parser {
	return s.sequence('[', ']')
}

func (s *state) dict() comb.Parser {
	return s.sequence('{', '}')
}

func (s *state) sequence(l, r rune) comb.Parser {
	return s.wrapChars(l, s.elems(), r)
}

func (s *state) comment() comb.Parser {
	return s.Void(s.And(s.Char(';'), s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *state) wrapChars(l rune, p comb.Parser, r rune) comb.Parser {
	return s.Wrap(s.And(s.Char(l), s.blank()), p, s.strippedChar(r))
}

func (s *state) strippedChar(r rune) comb.Parser {
	return s.strip(s.Char(r))
}

func (s *state) strip(p comb.Parser) comb.Parser {
	b := s.blank()
	return s.Wrap(b, p, b)
}

func (s *state) blank() comb.Parser {
	return s.Void(s.Many(s.Or(s.space(), s.comment())))
}

func (s *state) space() comb.Parser {
	return s.Void(s.Many1(s.InString(SPACE_CHARS)))
}

func (s *state) quote(p comb.Parser) comb.Parser {
	return s.And(s.Char('\''), p)
}

func (s *state) quotes(ps ...comb.Parser) []comb.Parser {
	qs := make([]comb.Parser, len(ps))

	for i, p := range ps {
		qs[i] = s.quote(p)
	}

	return qs
}

func (s *state) stringify(p comb.Parser) comb.Parser {
	f := func(any interface{}) interface{} {
		xs := any.([]interface{})
		rs := make([]rune, len(xs))

		for i, x := range xs {
			rs[i] = x.(rune)
		}

		return types.NewString(string(rs))
	}

	return s.App(f, p)
}
