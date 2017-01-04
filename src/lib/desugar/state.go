package desugar

type state struct {
	env map[string]interface{}
}

func newState() *state {
	return &state{make(map[string]interface{})}
}
