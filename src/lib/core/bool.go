package core

// BoolType represents a boolean values in the language.
type BoolType bool

// True is a true value.
var True = Normal(BoolType(true))

// False is a false value.
var False = Normal(BoolType(false))

// NewBool converts a Go boolean value into BoolType.
func NewBool(b bool) *Thunk {
	if b {
		return True
	}

	return False
}

// If returns the second argument when the first one is true or the third one
// otherwise.
var If = NewLazyFunction(
	NewSignature(nil, nil, "args", nil, nil, ""),
	func(ts ...*Thunk) Value {
		t := ts[0]

		for {
			v := t.Eval()
			l, ok := v.(ListType)

			if !ok {
				return NotListError(v)
			} else if l == emptyList {
				return argumentError("Not enough arguments to if function.")
			}

			v = l.rest.Eval()
			ll, ok := v.(ListType)

			if !ok {
				return NotListError(v)
			} else if ll == emptyList {
				return l.first
			}

			v = l.first.Eval()
			b, ok := v.(BoolType)

			if !ok {
				return NotBoolError(v)
			} else if b {
				return ll.first
			}

			t = ll.rest
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
