package core

type halfSignature struct {
	requireds []string
	optionals []OptionalArgument
	rest      string
}

func (as halfSignature) size() int {
	n := len(as.requireds) + len(as.optionals)

	if as.rest != "" {
		n++
	}

	return n
}
