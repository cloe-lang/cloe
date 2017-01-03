package parse

import "./comb"

type state struct{ *comb.State }

func newState(source string) *state {
	return &state{comb.NewState(source)}
}
