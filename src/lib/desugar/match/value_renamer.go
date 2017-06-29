package match

type valueRenamer struct {
	nameMap map[string]string
}

func newValueRenamer(m map[string]string) valueRenamer {
	return valueRenamer{m}
}

func (r valueRenamer) rename(x interface{}) interface{} {
	panic("Not implemented")
}
