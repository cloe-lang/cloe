package comb

import (
	"fmt"
)

// Block creates a parser parsing things which have the same indent.
func (s *State) Block(n int, p, q Parser) Parser {
	return func() (interface{}, error) {
		pos := s.linePosition
		p()
		return s.Many(s.atPosition(pos+n, q))()
	}
}

func (s *State) atPosition(pos int, p Parser) Parser {
	return func() (interface{}, error) {
		if s.linePosition != pos {
			return nil, fmt.Errorf("Invalid indent")
		}

		return p()
	}
}
