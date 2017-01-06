package vm

type String string

func NewString(s string) *Thunk {
	return Normal(String(s))
}
