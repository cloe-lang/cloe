package comb

import "strings"

// State represents a parser state.
type State struct {
	source               []rune
	lineNumber, position int
}

// NewState creates a parser state.
func NewState(source string) *State {
	return &State{source: ([]rune)(source)}
}

func (s State) exhausted() bool {
	return s.position >= len(s.source)
}

func (s State) currentRune() rune {
	if s.exhausted() {
		return '\x00'
	}

	return s.source[s.position]
}

func (s *State) increment() {
	if s.currentRune() == '\n' {
		s.lineNumber++
	}

	s.position++
}

// LineNumber returns a current line number.
func (s *State) LineNumber() int {
	return s.lineNumber
}

// Line returns a current line.
func (s *State) Line() string {
	return strings.Split(string(s.source), "\n")[s.lineNumber]
}
