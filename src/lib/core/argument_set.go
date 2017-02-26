package core

type argumentSet struct {
	requireds []string
	optionals []OptionalArgument
	rest      string
}

func (as argumentSet) size() int {
	n := len(as.requireds) + len(as.optionals)

	if as.rest != "" {
		n++
	}

	return n
}
