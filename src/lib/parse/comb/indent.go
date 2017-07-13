package comb

import (
	"fmt"
)

// Block creates a parser parsing things which have the same indent.
func (s *State) Block(n int, p, q Parser) Parser {
	return s.WithPosition(func() (interface{}, error) {
		x, err := p()
		if err != nil {
			return nil, err
		}

		xs, err := s.Many(s.atLinePosition(s.reference.linePosition+n, q))()
		if err != nil {
			return nil, err
		}

		return append([]interface{}{x}, xs.([]interface{})...), nil
	})
}

// WithPosition saves a current position and runs a given parser.
func (s *State) WithPosition(p Parser) Parser {
	return func() (interface{}, error) {
		s.savePosition()
		return p()
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
