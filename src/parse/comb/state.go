package comb

type State struct {
	source         []rune
	line, position uint
}

func NewState(source string) *State {
	return &State{source: ([]rune)(source)}
}

func (s State) currentRune() rune {
	return s.source[s.position]
}

func (s *State) increment() {
	if s.currentRune() == '\n' {
		s.line++
	}

	s.position++
}
