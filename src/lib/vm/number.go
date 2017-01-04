package vm

type Number float64

func NewNumber(n float64) Number {
	return Number(n)
}
