package vm

type String string

func NewString(s string) *Thunk {
	return Normal(String(s))
}

func (s String) Equal(e Equalable) bool {
	return s == e.(String)
}
