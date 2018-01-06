package core

import "math"

// NumberType represents a number in the language.
// It will perhaps be represented by DEC64 in the future release.
type NumberType float64

// NewNumber creates a thunk containing a number value.
func NewNumber(n float64) *Thunk {
	return Normal(NumberType(n))
}

// Add sums up numbers of arguments.
var Add = newCommutativeOperator(0, func(n, m NumberType) NumberType { return n + m })

// Sub subtracts arguments of the second to the last from the first one as numbers.
var Sub = newInverseOperator(func(n, m NumberType) NumberType { return n - m })

// Mul multiplies numbers of arguments.
var Mul = newCommutativeOperator(1, func(n, m NumberType) NumberType { return n * m })

// Div divides the first argument by arguments of the second to the last one by one.
var Div = newInverseOperator(func(n, m NumberType) NumberType { return n / m })

// FloorDiv divides the first argument by arguments of the second to the last one by one.
var FloorDiv = newInverseOperator(func(n, m NumberType) NumberType {
	return NumberType(math.Floor(float64(n / m)))
})

// Mod calculate a remainder of a division of the first argument by the second one.
var Mod = newBinaryOperator(math.Mod)

// Pow calculates an exponentiation from a base of the first argument and an
// exponent of the second argument.
var Pow = newBinaryOperator(math.Pow)

func newCommutativeOperator(i NumberType, f func(n, m NumberType) NumberType) *Thunk {
	return NewLazyFunction(
		NewSignature(nil, nil, "nums", nil, nil, ""),
		func(ts ...*Thunk) Value {
			l, err := ts[0].EvalList()

			if err != nil {
				return err
			}

			ts, e := l.ToValues()

			if e != nil {
				return e
			}

			a := i

			for _, t := range ts {
				n, err := t.EvalNumber()

				if err != nil {
					return err
				}

				a = f(a, n)
			}

			return a
		})
}

func newInverseOperator(f func(n, m NumberType) NumberType) *Thunk {
	return NewLazyFunction(
		NewSignature([]string{"initial"}, nil, "nums", nil, nil, ""),
		func(ts ...*Thunk) Value {
			a, err := ts[0].EvalNumber()

			if err != nil {
				return err
			}

			l, err := ts[1].EvalList()

			if err != nil {
				return err
			}

			ts, e := l.ToValues()

			if e != nil {
				return e
			}

			for _, t := range ts {
				n, err := t.EvalNumber()

				if err != nil {
					return err
				}

				a = f(a, n)
			}

			return a
		})
}

func newBinaryOperator(f func(n, m float64) float64) *Thunk {
	return NewStrictFunction(
		NewSignature([]string{"first", "second"}, nil, "", nil, nil, ""),
		func(ts ...*Thunk) Value {
			ns := [2]NumberType{}

			for i, t := range ts {
				n, err := t.EvalNumber()

				if err != nil {
					return err
				}

				ns[i] = n
			}

			return NumberType(f(float64(ns[0]), float64(ns[1])))
		})
}

var isInt = NewLazyFunction(
	NewSignature([]string{"number"}, nil, "", nil, nil, ""),
	func(ts ...*Thunk) Value {
		n, err := ts[0].EvalNumber()

		if err != nil {
			return err
		}

		return NewBool(math.Mod(float64(n), 1) == 0)
	})

func (n NumberType) compare(c comparable) int {
	if n < c.(NumberType) {
		return -1
	} else if n > c.(NumberType) {
		return 1
	}

	return 0
}

func (NumberType) ordered() {}

func (n NumberType) string() Value {
	return sprint(n)
}
