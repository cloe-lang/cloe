package comb

import "strings"

type State struct {
	source               []rune
	lineNumber, position int
}

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

func (s *State) LineNumber() int {
	return s.lineNumber
}

func (s *State) Line() string {
	return strings.Split(string(s.source), "\n")[s.lineNumber]
}
