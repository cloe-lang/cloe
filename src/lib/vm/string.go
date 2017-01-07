package vm

type String string

func NewString(s string) *Thunk {
	return Normal(String(s))
}

func (s String) Equal(o Object) bool {
	return s == o.(String)
}
