package vm

type Number float64

func NewNumber(n float64) *Thunk {
	return Normal(Number(n))
}

func (n Number) Equal(e Equalable) bool {
	return n == e.(Number)
}
