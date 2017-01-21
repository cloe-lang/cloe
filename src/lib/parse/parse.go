package parse

import "./comb"

const (
	bracketChars = "()[]{}"
	commentChar  = ';'
	invalidChars = "\x00"
	quoteChar    = '`'
	quoteString  = "quote"
	spaceChars   = " \t\n\r"
	specialChar  = '$'
)

func Parse(source string) []interface{} {
	m, err := newState(source).module()()

	if err != nil {
		panic(err.Error())
	}

	return m.([]interface{})
}

func (s *state) module() comb.Parser {
	return s.Exhaust(s.Wrap(s.blank(), s.expressions(), s.None()))
}

func (s *state) expressions() comb.Parser {
	return s.Lazy(s.strictExpressions)
}

func (s *state) strictExpressions() comb.Parser {
	return s.Many(s.expression())
}

func (s *state) expression() comb.Parser {
	ps := []comb.Parser{s.atom(), s.list(), s.array(), s.dict()}

	return s.strip(s.Or(append(ps, s.quotes(ps...)...)...))
}

func (s *state) atom() comb.Parser {
	return s.Or(s.stringLiteral(), s.identifier())
}

func (s *state) identifier() comb.Parser {
	return s.stringify(s.Many1(s.NotInString(
		bracketChars + string(commentChar) + invalidChars + string(quoteChar) +
			spaceChars + string(specialChar))))
}

func (s *state) stringLiteral() comb.Parser {
	b := s.blank()
	c := s.Char('"')

	return s.App(prependQuote, s.stringify(s.Wrap(
		c,
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		s.And(c, b))))
}

func prependQuote(x interface{}) interface{} {
	return []interface{}{quoteString, x}
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
	return s.wrapChars(l, s.expressions(), r)
}

func (s *state) comment() comb.Parser {
	return s.Void(s.And(
		s.Char(commentChar),
		s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *state) wrapChars(l rune, p comb.Parser, r rune) comb.Parser {
	return s.Wrap(s.strippedChar(l), p, s.strippedChar(r))
}

func (s *state) strippedChar(r rune) comb.Parser {
	return s.strip(s.Char(r))
}

func (s *state) strip(p comb.Parser) comb.Parser {
	b := s.blank()
	return s.Wrap(s.None(), p, b)
}

func (s *state) blank() comb.Parser {
	return s.Void(s.Many(s.Or(s.space(), s.comment())))
}

func (s *state) space() comb.Parser {
	return s.Void(s.Many1(s.InString(spaceChars)))
}

func (s *state) quote(p comb.Parser) comb.Parser {
	return s.And(s.Replace(quoteString, s.Char(quoteChar)), p)
}

func (s *state) quotes(ps ...comb.Parser) []comb.Parser {
	qs := make([]comb.Parser, len(ps))

	for i, p := range ps {
		qs[i] = s.quote(p)
	}

	return qs
}

func (s *state) stringify(p comb.Parser) comb.Parser {
	return s.App(stringifyChars, p)
}

func stringifyChars(any interface{}) interface{} {
	xs := any.([]interface{})
	rs := make([]rune, len(xs))

	for i, x := range xs {
		rs[i] = x.(rune)
	}

	return string(rs)
}
