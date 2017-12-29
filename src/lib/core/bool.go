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
	NewSignature(nil, nil, "conds", nil, nil, ""),
	func(ts ...*Thunk) Value {
		v := ts[0].Eval()
		l, ok := v.(ListType)

		if !ok {
			return NotListError(v)
		}

		ts, err := l.ToThunks()

		if err != nil {
			return err
		}

		if len(ts)%2 == 0 {
			return argumentError("Number of arguments of if function must be even but %v.", len(ts))
		}

		for i := 0; i < len(ts)-2; i += 2 {
			v := ts[i].Eval()
			b, ok := v.(BoolType)

			if !ok {
				return NotBoolError(v)
			}

			if b {
				return ts[i+1]
			}
		}

		return ts[len(ts)-1]
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
