package comb

import (
	"fmt"
	"strings"
)

// Char creates a parser parsing a character.
func (s *State) Char(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() != r {
			return nil, fmt.Errorf("invalid character, '%c'", s.currentRune())
		}

		s.increment()

		return r, nil
	}
}

// NotChar creates a parser parsing a character which is not one of an argument.
func (s *State) NotChar(r rune) Parser {
	return func() (interface{}, error) {
		if s.currentRune() == r {
			return nil, fmt.Errorf("should not be '%c'", r)
		}

		defer s.increment()

		return s.currentRune(), nil
	}
}

// String creates a parser parsing a string.
func (s *State) String(str string) Parser {
	rs := ([]rune)(str)
	ps := make([]Parser, len(rs))

	for i, r := range rs {
		ps[i] = s.Char(r)
	}

	return s.Stringify(s.And(ps...))
}

// InString creates a parser parsing a character in a given string.
func (s *State) InString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; ok {
			defer s.increment()
			return s.currentRune(), nil
		}

		return nil, fmt.Errorf("invalid character, '%c'", s.currentRune())
	}
}

// NotInString creates a parser parsing a character not in a given string.
func (s *State) NotInString(str string) Parser {
	rs := stringToRuneSet(str)

	return func() (interface{}, error) {
		if _, ok := rs[s.currentRune()]; !ok {
			defer s.increment()
			return s.currentRune(), nil
		}

		return nil, fmt.Errorf("invalid character, '%c'", s.currentRune())
	}
}

// Wrap wraps a parser with parsers which parse something on the leftside and
// rightside of it and creates a new parser. Its parsing result will be m's.
func (s *State) Wrap(l, m, r Parser) Parser {
	return second(s.And(l, m, r))
}

// Prefix creates a parser with a prefix parser and content parser and returns
// the latter's result.
func (s *State) Prefix(pre, p Parser) Parser {
	return second(s.And(pre, p))
}

func second(p Parser) Parser {
	return func() (interface{}, error) {
		results, err := p()

		if results, ok := results.([]interface{}); ok {
			return results[1], err
		}

		return nil, err
	}
}

// Many creates a parser of more than or equal to 0 reptation of a given parser.
func (s *State) Many(p Parser) Parser {
	return func() (interface{}, error) {
		results, err := s.Many1(p)()

		if err != nil {
			return []interface{}{}, nil
		}

		return results, nil
	}
}

// Many1 creates a parser of more than 0 reptation of a given parser.
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

// Or creates a selectional parser from given parsers.
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

// And creates a parser combining given parsers sequentially.
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

// Lazy evaluates and runs a given parser constructor. This is useful to define
// recursive parsers.
func (s *State) Lazy(f func() Parser) Parser {
	p := Parser(nil)

	return func() (interface{}, error) {
		if p == nil {
			p = f()
		}

		return p()
	}
}

// Void creates a parser whose result is always nil from a given parser.
func (State) Void(p Parser) Parser {
	return func() (interface{}, error) {
		_, err := p()
		return nil, err
	}
}

// Exhaust creates a parser which fails when a source string is not exhausted
// after running a given parser.
func (s *State) Exhaust(p Parser) Parser {
	return func() (interface{}, error) {
		result, err := p()

		if err != nil {
			return result, err
		} else if !s.exhausted() {
			return nil, fmt.Errorf(
				"Some characters are left in source. %#v",
				string(s.source[s.current.sourcePosition:]))
		}

		return result, err
	}
}

// App applies a function to results of a given parser.
func (s *State) App(f func(interface{}) interface{}, p Parser) Parser {
	return func() (interface{}, error) {
		result, err := p()

		if err == nil {
			return f(result), err
		}

		return result, err
	}
}

// Replace replaces a result of a given parser and creates a new parser.
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

// None creates a parser which parses nothing and succeeds always.
func (s *State) None() Parser {
	return func() (interface{}, error) {
		return nil, nil
	}
}

// Maybe creates a parser which runs a given parser or parses nothing when it
// fails.
func (s *State) Maybe(p Parser) Parser {
	return s.Or(p, s.None())
}

// Stringify creates a parser which returns a string converted from a result of
// a given parser. The result of a given parser must be a rune, a string or a
// sequence of them in []interface{}.
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

	panic("Unreachable")
}
