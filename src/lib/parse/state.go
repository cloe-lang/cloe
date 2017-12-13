package parse

import (
	"github.com/coel-lang/coel/src/lib/debug"
	"github.com/coel-lang/coel/src/lib/parse/comb"
)

type state struct {
	*comb.State
	file string
}

func newState(file, source string) *state {
	return &state{comb.NewState(source), file}
}

func (s state) debugInfo() debug.Info {
	return debug.NewInfo(s.file, s.LineNumber(), s.LinePosition(), s.Line())
}
