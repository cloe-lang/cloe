package core

import (
	"hash/fnv"

	"github.com/raviqqe/hamt"
)

// BooleanType represents a boolean values in the language.
type BooleanType bool

// Eval evaluates a value into a WHNF.
func (b *BooleanType) eval() Value {
	return b
}

var (
	trueStruct, falseStruct = BooleanType(true), BooleanType(false)

	// True is a true value.
	True = &trueStruct

	// False is a false value.
	False = &falseStruct
)

// NewBoolean converts a Go boolean value into BooleanType.
func NewBoolean(b bool) *BooleanType {
	if b {
		return True
	}

	return False
}

// If returns the second argument when the first one is true or the third one
// otherwise.
var If = NewLazyFunction(
	NewSignature(nil, "args", nil, ""),
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

			b, err := EvalBoolean(l.First())

			if err != nil {
				return err
			} else if b {
				return ll.First()
			}

			v = ll.Rest()
		}
	})

func (b *BooleanType) compare(c comparable) int {
	if *b == *c.(*BooleanType) {
		return 0
	} else if *b {
		return 1
	}

	return -1
}

// Hash hashes a value.
func (b *BooleanType) Hash() uint32 {
	h := fnv.New32()

	if *b {
		if _, err := h.Write([]byte{1}); err != nil {
			panic(err)
		}
	} else {
		if _, err := h.Write([]byte{0}); err != nil {
			panic(err)
		}
	}

	return h.Sum32()
}

// Equal checks equality.
func (b *BooleanType) Equal(e hamt.Entry) bool {
	if c, ok := e.(*BooleanType); ok {
		return *b == *c
	}

	return false
}

func (b *BooleanType) string() Value {
	return sprint(*b)
}
