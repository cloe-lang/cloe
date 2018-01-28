package core

import "math"

// NumberType represents a number in the language.
// It will perhaps be represented by DEC64 in the future release.
type NumberType float64

// Eval evaluates a value into a WHNF.
func (n *NumberType) eval() Value {
	return n
}

// NewNumber creates a thunk containing a number value.
func NewNumber(n float64) *NumberType {
	m := NumberType(n)
	return &m
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

func newCommutativeOperator(i NumberType, f func(n, m NumberType) NumberType) Value {
	return NewLazyFunction(
		NewSignature(nil, nil, "nums", nil, nil, ""),
		func(vs ...Value) Value {
			l, err := EvalList(vs[0])

			if err != nil {
				return err
			}

			a := i

			for !l.Empty() {
				n, err := EvalNumber(l.First())

				if err != nil {
					return err
				}

				a = f(a, n)

				l, err = EvalList(l.Rest())

				if err != nil {
					return err
				}
			}

			return &a
		})
}

func newInverseOperator(f func(n, m NumberType) NumberType) Value {
	return NewLazyFunction(
		NewSignature([]string{"initial"}, nil, "nums", nil, nil, ""),
		func(vs ...Value) Value {
			a, err := EvalNumber(vs[0])

			if err != nil {
				return err
			}

			l, err := EvalList(vs[1])

			if err != nil {
				return err
			}

			for !l.Empty() {
				n, err := EvalNumber(l.First())

				if err != nil {
					return err
				}

				a = f(a, n)

				l, err = EvalList(l.Rest())

				if err != nil {
					return err
				}
			}

			return &a
		})
}

func newBinaryOperator(f func(n, m float64) float64) Value {
	return NewStrictFunction(
		NewSignature([]string{"first", "second"}, nil, "", nil, nil, ""),
		func(vs ...Value) Value {
			ns := [2]NumberType{}

			for i, t := range vs {
				n, err := EvalNumber(t)

				if err != nil {
					return err
				}

				ns[i] = n
			}

			return NewNumber(f(float64(ns[0]), float64(ns[1])))
		})
}

// IsInt checks if a number value is an integer or not.
func IsInt(n NumberType) bool {
	return math.Mod(float64(n), 1) == 0
}

func (n *NumberType) compare(c comparable) int {
	if *n < *c.(*NumberType) {
		return -1
	} else if *n > *c.(*NumberType) {
		return 1
	}

	return 0
}

func (*NumberType) ordered() {}

func (n *NumberType) string() Value {
	return sprint(*n)
}
