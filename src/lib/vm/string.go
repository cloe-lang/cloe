package vm

type String string

func NewString(s string) String {
	return String(s)
}

func StringThunk(s string) *Thunk {
	return Normal(NewString(s))
}

func (s String) Equal(e Equalable) Object {
	return rawBool(s == e)
}

func (s String) Add(a Addable) Addable {
	return s + a.(String)
}
