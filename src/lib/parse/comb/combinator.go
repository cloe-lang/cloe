package comb

import (
	"fmt"
	"strings"
)

func (s *State) Char(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() != r {
			return nil, NewError("Invalid character. '%c'", s.currentRune())
		}

		s.increment()

		return r, nil
	}
}

func (s *State) NotChar(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() == r {
			return nil, NewError("Should not be '%c'.", r)
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

	return s.Stringify(s.And(ps...))
}

func (s *State) InString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; ok {
			defer s.increment()
			return s.currentRune(), nil
		}

		return nil, NewError("Invalid character. '%c'", s.currentRune())
	}
}

func (s *State) NotInString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; !ok {
			defer s.increment()
			return s.currentRune(), nil
		}

		return nil, NewError("Invalid character. '%c'", s.currentRune())
	}
}

func (s *State) Wrap(l, m, r Parser) Parser {
	p := s.And(l, m, r)

	return func() (interface{}, error) {
		results, err := p()

		if results, ok := results.([]interface{}); ok {
			return results[1], err
		}

		return results, err
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

func (s *State) Or(ps ...Parser) Parser {
	return func() (interface{}, error) {
		var err error
		old := *s

		for _, p := range ps {
			var result interface{}

			result, err = p()

			if err == nil {
				return result, nil
			}

			*s = old
		}

		return nil, err
	}
}

func (s *State) And(ps ...Parser) Parser {
	return func() (interface{}, error) {
		results := make([]interface{}, len(ps))
		old := *s

		for i, p := range ps {
			result, err := p()

			if err != nil {
				*s = old
				return nil, err
			}

			results[i] = result
		}

		return results, nil
	}
}

func (s *State) Lazy(f func() Parser) Parser {
	p := Parser(nil)

	return func() (interface{}, error) {
		if p == nil {
			p = f()
		}

		return p()
	}
}

func (State) Void(p Parser) Parser {
	return func() (interface{}, error) {
		_, err := p()
		return nil, err
	}
}

func (s *State) Exhaust(p Parser) Parser {
	return func() (interface{}, error) {
		result, err := p()

		if !s.exhausted() {
			return nil, NewError(
				"Some characters are left in source. %#v",
				string(s.source[s.position:]))
		}

		return result, err
	}
}

func (s *State) App(f func(interface{}) interface{}, p Parser) Parser {
	return func() (interface{}, error) {
		result, err := p()

		if err == nil {
			return f(result), err
		}

		return result, err
	}
}

func (s *State) Replace(x interface{}, p Parser) Parser {
	return s.App(func(_ interface{}) interface{} { return x }, p)
}

func stringToRuneSet(s string) map[rune]bool {
	rs := make(map[rune]bool)

	for _, r := range s {
		rs[r] = true
	}

	return rs
}

func (s *State) None() Parser {
	return func() (interface{}, error) {
		return nil, nil
	}
}

func (s *State) Maybe(p Parser) Parser {
	return s.Or(p, s.None())
}

func (s *State) Stringify(p Parser) Parser {
	return s.App(func(x interface{}) interface{} { return stringify(x) }, p)
}

func stringify(x interface{}) string {
	switch x := x.(type) {
	case string:
		return x
	case rune:
		return string(x)
	case []interface{}:
		ss := make([]string, len(x))

		for i, s := range x {
			ss[i] = stringify(s)
		}

		return strings.Join(ss, "")
	}

	panic(fmt.Sprint("Invalid type.", x))
}
