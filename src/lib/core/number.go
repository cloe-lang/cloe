package core

import "math"

// NumberType represents a number in the language.
// It will perhaps be represented by DEC64 in the future release.
type NumberType float64

// NewNumber creates a thunk containing a number value.
func NewNumber(n float64) *Thunk {
	return Normal(NumberType(n))
}

func (n NumberType) equal(e equalable) Value {
	return rawBool(n == e.(NumberType))
}

// Add sums up numbers of arguments.
var Add = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "nums",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		vs, err := l.ToValues()

		if err != nil {
			return err
		}

		sum := NumberType(0)

		for _, v := range vs {
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			sum += n
		}

		return sum
	})

// Sub subtracts arguments of the second to the last from the first one as numbers.
var Sub = NewLazyFunction(
	NewSignature(
		[]string{"minuend"}, []OptionalArgument{}, "subtrahends",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n0, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		vs, err := l.ToValues()

		if err != nil {
			return err
		}

		if len(vs) == 0 {
			return NumArgsError("sub", ">= 1")
		}

		for _, v := range vs {
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			n0 -= n
		}

		return n0
	})

// Mul multiplies numbers of arguments.
var Mul = NewLazyFunction(
	NewSignature(
		[]string{}, []OptionalArgument{}, "nums",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		vs, err := l.ToValues()

		if err != nil {
			return err
		}

		prod := NumberType(1)

		for _, v := range vs {
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			prod *= n
		}

		return prod
	})

// Div divides the first argument by arguments of the second to the last one by one.
var Div = NewLazyFunction(
	NewSignature(
		[]string{"dividend"}, []OptionalArgument{}, "divisors",
		[]string{}, []OptionalArgument{}, "",
	),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		n0, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = ts[1].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		vs, err := l.ToValues()

		if err != nil {
			return err
		}

		if len(vs) == 0 {
			return NumArgsError("div", ">= 1")
		}

		for _, v := range vs {
			n, ok := v.(NumberType)

			if !ok {
				return NotNumberError(v)
			}

			n0 /= n
		}

		return n0
	})

// TODO: Implement FloorDiv function.

// Mod calculate a remainder of a division of the first argument by the second one.
var Mod = NewStrictFunction(
	NewSignature(
		[]string{"dividend", "divisor"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		if len(vs) != 2 {
			return NumArgsError("mod", "2")
		}

		v := vs[0]
		n1, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = vs[1]
		n2, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		return NewNumber(math.Mod(float64(n1), float64(n2)))
	})

// Pow calculates an exponentiation from a base of the first argument and an
// exponent of the second argument.
var Pow = NewStrictFunction(
	NewSignature(
		[]string{"base", "exponent"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		n1, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		v = vs[1]
		n2, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		return NewNumber(math.Pow(float64(n1), float64(n2)))
	})

var isInt = NewStrictFunction(
	NewSignature(
		[]string{"number"}, []OptionalArgument{}, "",
		[]string{}, []OptionalArgument{}, "",
	),
	func(vs ...Value) Value {
		v := vs[0]
		n, ok := v.(NumberType)

		if !ok {
			return NotNumberError(v)
		}

		return NewBool(math.Mod(float64(n), 1) == 0)
	})

func (n NumberType) less(o ordered) bool {
	return n < o.(NumberType)
}

func (n NumberType) string() Value {
	return sprint(n)
}
