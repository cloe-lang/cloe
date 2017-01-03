package parse

const SPECIAL_CHARS = "()[]',"

type Parser func() (interface{}, error)

func (s *State) Module() Parser {
	return s.Elems()
}

func (s *State) Elems() Parser {
	return s.Many(s.Elem())
}

func (s *State) Elem() Parser {
	return s.Or(s.Atom(), s.List())
}

func (s *State) Atom() Parser {
	return s.Or(s.Identifier(), s.StringLiteral())
}

func (s *State) Identifier() Parser {
	return s.Many1(s.NotInString(SPECIAL_CHARS))
}

func (s *State) StringLiteral() Parser {
	return s.Wrap(
		'"',
		s.Many(s.Or(s.NotInString("\"\\"), s.String("\\\""), s.String("\\\\"))),
		'"')
}

func (s *State) List() Parser {
	return s.Wrap('(', s.Elems(), ')')
}

func (s *State) Wrap(l rune, m Parser, r rune) Parser {
	return func() (interface{}, error) {
		s.StrippedChar(l)()
		result, err := m()
		s.StrippedChar(r)()
		return result, err
	}
}

func (s *State) StrippedChar(r rune) Parser {
	return s.Strip(s.Char('('))
}

func (s *State) Char(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() != r {
			return nil, NewError("Invalid character. ('%c')", s.currentRune())
		}

		s.Increment()

		return r, nil
	}
}

func (s *State) NotChar(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() == r {
			return nil, NewError("Should not be %c.", r)
		}

		defer s.Increment()

		return s.currentRune(), nil
	}
}

func (s *State) String(str string) Parser {
	rs := ([]rune)(str)
	ps := make([]Parser, len(rs))

	for i, r := range rs {
		ps[i] = s.Char(r)
	}

	return s.And(ps...)
}

func (s *State) InString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; ok {
			defer s.Increment()
			return s.currentRune(), nil
		}

		return nil, NewError("Invalid character. (%c)", s.currentRune())
	}
}

func (s *State) NotInString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; !ok {
			defer s.Increment()
			return s.currentRune(), nil
		}

		return nil, NewError("Invalid character. (%c)", s.currentRune())
	}
}

func (s *State) Comment() Parser {
	return s.Void(s.And(s.Char(';'), s.Many(s.NotChar('\n')), s.Char('\n')))
}

func (s *State) Strip(p Parser) Parser {
	return func() (interface{}, error) {
		s.Blank()()
		result, err := p()
		s.Blank()()
		return result, err
	}
}

func (s *State) Blank() Parser {
	return s.Void(s.Many(s.Or(s.Space(), s.Comment())))
}

func (s *State) Space() Parser {
	return func() (interface{}, error) {
		_, err := s.Many1(s.Or(
			s.Char(' '),
			s.Char(','),
			s.Char('\t'),
			s.Char('\n'),
			s.Char('\r'),
		))()

		return nil, err
	}
}

func (s *State) Many(p Parser) Parser {
	return func() (interface{}, error) {
		results, err := s.Many1(p)()

		if err != nil {
			return []interface{}{}, nil
		}

		return results, nil
	}
}

func (s *State) Many1(p Parser) Parser {
	return func() (interface{}, error) {
		var results []interface{}

		for i := 0; ; i++ {
			result, err := p()

			if err != nil && i == 0 {
				return nil, err
			} else if err != nil {
				break
			}

			results = append(results, result)
		}

		return results, nil
	}
}

func (State) Or(ps ...Parser) Parser {
	return func() (interface{}, error) {
		var err error

		for _, p := range ps {
			var result interface{}

			result, err = p()

			if err == nil {
				return result, nil
			}
		}

		return nil, err
	}
}

func (State) And(ps ...Parser) Parser {
	return func() (interface{}, error) {
		results := make([]interface{}, len(ps))

		for i, p := range ps {
			result, err := p()

			if err != nil {
				return nil, err
			}

			results[i] = result
		}

		return results, nil
	}
}

func (State) Void(p Parser) Parser {
	return func() (interface{}, error) {
		_, err := p()
		return nil, err
	}
}

func stringToRuneSet(s string) map[rune]bool {
	var rs map[rune]bool

	for _, r := range s {
		rs[r] = true
	}

	return rs
}
