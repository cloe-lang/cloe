package core

// BoolType represents a boolean values in the language.
type BoolType bool

// Eval evaluates a value into a WHNF.
func (b BoolType) eval() Value {
	return b
}

// True is a true value.
var True = BoolType(true)

// False is a false value.
var False = BoolType(false)

// NewBool converts a Go boolean value into BoolType.
func NewBool(b bool) BoolType {
	return BoolType(b)
}

// If returns the second argument when the first one is true or the third one
// otherwise.
var If = NewLazyFunction(
	NewSignature(nil, nil, "args", nil, nil, ""),
	func(vs ...Value) Value {
		v := vs[0]

		for {
			l, err := EvalList(v)

			if err != nil {
				return err
			}

			ll, err := EvalList(l.Rest())

			if err != nil {
				return err
			} else if ll.Empty() {
				return l.First()
			}

			b, err := EvalBool(l.First())

			if err != nil {
				return err
			} else if b {
				return ll.First()
			}

			v = ll.Rest()
		}
	})

func (b BoolType) compare(c comparable) int {
	if b == c.(BoolType) {
		return 0
	} else if b {
		return 1
	}

	return -1
}

func (b BoolType) string() Value {
	return sprint(b)
}
