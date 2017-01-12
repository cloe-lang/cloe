package vm

type stringType string

func NewString(s string) *Thunk {
	return Normal(stringType(s))
}

func (s stringType) equal(e equalable) Object {
	return rawBool(s == e)
}

func (s stringType) Add(a Addable) Addable {
	return s + a.(stringType)
}
