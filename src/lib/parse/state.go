package parse

import "github.com/raviqqe/tisp/src/lib/parse/comb"

type state struct{ *comb.State }

func newState(source string) *state {
	return &state{comb.NewState(source)}
}
