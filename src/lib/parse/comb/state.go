package comb

import "strings"

type position struct {
	lineNumber, linePosition, sourcePosition int
}

// State represents a parser state.
type State struct {
	source             []rune
	current, reference position
}

// NewState creates a parser state.
func NewState(source string) *State {
	return &State{source: ([]rune)(source)}
}

func (s State) exhausted() bool {
	return s.current.sourcePosition >= len(s.source)
}

func (s State) currentRune() rune {
	if s.exhausted() {
		return '\x00'
	}

	return s.source[s.current.sourcePosition]
}

func (s *State) increment() {
	if s.currentRune() == '\n' {
		s.current.lineNumber++
		s.current.linePosition = 0
	} else {
		s.current.linePosition++
	}

	s.current.sourcePosition++
}

// LineNumber returns a current line number.
func (s *State) LineNumber() int {
	return s.current.lineNumber
}

// LinePosition returns a current line number.
func (s *State) LinePosition() int {
	return s.current.linePosition
}

// Line returns a current line.
func (s *State) Line() string {
	return strings.Split(string(s.source), "\n")[s.current.lineNumber]
}
