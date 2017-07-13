package comb

import (
	"fmt"
)

// Block creates a parser parsing things which have the same indent.
func (s *State) Block(n int, p, q, r Parser) Parser {
	return s.WithPosition(func() (interface{}, error) {
		x, err := p()
		if err != nil {
			return nil, err
		}

		ys, err := s.Many(s.atLinePosition(s.reference.linePosition+n, q))()
		if err != nil {
			return nil, err
		}

		z, err := r()
		if err != nil {
			return nil, err
		}

		return append([]interface{}{x}, append(ys.([]interface{}), z)...), nil
	})
}

// SameLine creates a parser parsing 2 things in the same line.
func (s *State) SameLine(p, q Parser) Parser {
	return s.WithPosition(func() (interface{}, error) {
		x, err := p()
		if err != nil {
			return nil, err
		}

		if s.current.lineNumber != s.reference.lineNumber {
			return nil, fmt.Errorf("Invalid new line is detected")
		}

		y, err := q()
		if err != nil {
			return nil, err
		}

		return []interface{}{x, y}, nil
	})
}

// SameLineOrIndented creates a parser whiich parses something in the same line
// or indented. It retrieves a reference position saved by WithPosition
// combinator implicitly differently from other indent-aware combinators.
func (s *State) SameLineOrIndented(n int, p Parser) Parser {
	return func() (interface{}, error) {
		if s.current.lineNumber != s.reference.lineNumber &&
			s.current.linePosition < s.reference.linePosition+n {
			return nil, fmt.Errorf("Invalid indent")
		}

		return p()
	}
}

// WithPosition saves a current position and runs a given parser.
func (s *State) WithPosition(p Parser) Parser {
	return func() (interface{}, error) {
		old := s.reference
		s.reference = s.current
		x, err := p()
		s.reference = old
		return x, err
	}
}

func (s *State) atLinePosition(pos int, p Parser) Parser {
	return func() (interface{}, error) {
		if s.current.linePosition != pos {
			return nil, fmt.Errorf("Invalid indent")
		}

		return p()
	}
}
