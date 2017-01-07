package vm

type Number float64

func NewNumber(n float64) *Thunk {
	return Normal(Number(n))
}

func (n Number) Equal(o Object) bool {
	return n == o.(Number)
}
