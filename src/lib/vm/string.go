package vm

type stringType string

func NewString(s string) *Thunk {
	return Normal(stringType(s))
}

func (s stringType) Equal(e Equalable) Object {
	return rawBool(s == e)
}

func (s stringType) Add(a Addable) Addable {
	return s + a.(stringType)
}
