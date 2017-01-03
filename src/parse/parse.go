package parse

import (
	"../types"
	"./comb"
)

func Parse(source string) types.Object {
	o, err := newState(source).module()()

	if err != nil {
		panic(err.Error())
	}

	return o
}

func (s *state) module() comb.Parser {
	return s.elems()
}

func (s *state) elems() comb.Parser {
	return s.Many(s.elem())
}

func (s *state) elem() comb.Parser {
	return s.Or(s.atom(), s.list())
}

func (s *state) atom() comb.Parser {
	return s.Or(s.identifier(), s.stringLiteral())
}

func (s *state) identifier() comb.Parser {
	return s.Many1(s.NotInString("()[]',"))
}

func (s *state) stringLiteral() comb.Parser {
	return s.wrap(
		'"',
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		'"')
}

func (s *state) list() comb.Parser {
	return s.wrap('(', s.elems(), ')')
}

func (s *state) comment() comb.Parser {
	return s.Void(s.And(s.Char(';'), s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *state) wrap(l rune, p comb.Parser, r rune) comb.Parser {
	return func() (interface{}, error) {
		s.strippedChar(l)()
		result, err := p()
		s.strippedChar(r)()
		return result, err
	}
}

func (s *state) strippedChar(r rune) comb.Parser {
	return s.strip(s.Char(r))
}

func (s *state) strip(p comb.Parser) comb.Parser {
	return func() (interface{}, error) {
		s.blank()()
		result, err := p()
		s.blank()()
		return result, err
	}
}

func (s *state) blank() comb.Parser {
	return s.Void(s.Many(s.Or(s.space(), s.comment())))
}

func (s *state) space() comb.Parser {
	return s.Void(s.Many1(s.InString(" ,\t\n\r")))
}
