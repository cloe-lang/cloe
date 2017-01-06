package vm

type Number float64

func NewNumber(n float64) *Thunk {
	return Normal(Number(n))
}
