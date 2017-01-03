package comb

func (s *State) Char(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() != r {
			return nil, NewError("Invalid character. ('%c')", s.currentRune())
		}

		s.increment()

		return r, nil
	}
}

func (s *State) NotChar(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() == r {
			return nil, NewError("Should not be %c.", r)
		}

		defer s.increment()

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
			defer s.increment()
			return s.currentRune(), nil
		}

		return nil, NewError("Invalid character. (%c)", s.currentRune())
	}
}

func (s *State) NotInString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; !ok {
			defer s.increment()
			return s.currentRune(), nil
		}

		return nil, NewError("Invalid character. (%c)", s.currentRune())
	}
}

func (s *State) Wrap(l, m, r Parser) Parser {
	return func() (interface{}, error) {
		l()
		result, err := m()
		r()
		return result, err
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
	rs := make(map[rune]bool)

	for _, r := range s {
		rs[r] = true
	}

	return rs
}
