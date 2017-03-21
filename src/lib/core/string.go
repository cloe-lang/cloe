package core

import (
	"strings"
)

// StringType represents strings in the language.
type StringType string

// NewString creates a string in the language from one in Go.
func NewString(s string) *Thunk {
	return Normal(StringType(s))
}

func (s StringType) equal(e equalable) Value {
	return rawBool(s == e)
}

func (s StringType) call(args Arguments) Value {
	return Index.Eval().(callable).call(NewPositionalArguments(Normal(s)).Merge(args))
}

func (s StringType) index(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	i := int(n)

	if i >= len(s) {
		return OutOfRangeError()
	}

	return s[i : i+1]
}

func (s StringType) insert(ts ...*Thunk) Value {
	v := ts[0].Eval()
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	i := int(n)

	if i > len(s) {
		return OutOfRangeError()
	}

	v = ts[1].Eval()
	ss, ok := v.(StringType)

	if !ok {
		return NotStringError(v)
	}

	return s[:i] + ss + s[i:]
}

func (s StringType) merge(ts ...*Thunk) Value {
	ts = append([]*Thunk{Normal(s)}, ts...)

	for _, t := range ts {
		go t.Eval()
	}

	ss := make([]string, len(ts))

	for i, t := range ts {
		v := t.Eval()
		s, ok := v.(StringType)

		if !ok {
			return TypeError(v, "String")
		}

		ss[i] = string(s)
	}

	return StringType(strings.Join(ss, ""))
}

func (s StringType) delete(v Value) Value {
	n, err := checkIndex(v)
	i := int(n)

	if err != nil {
		return err
	} else if i >= len(s) {
		return OutOfRangeError()
	}

	return s[:i] + s[i+1:]
}

func (s StringType) toList() Value {
	if s == "" {
		return emptyList
	}

	rs := []rune(string(s))

	return cons(
		NewString(string(rs[0])),
		PApp(ToList, NewString(string(rs[1:]))))
}

func (s StringType) less(o ordered) bool {
	return s < o.(StringType)
}

func (s StringType) string() Value {
	return s
}

func (s StringType) size() Value {
	return NumberType(len(([]rune)(string(s))))
}

func (s StringType) include(v Value) Value {
	ss, ok := v.(StringType)

	if !ok {
		return NotStringError(v)
	}

	return NewBool(strings.Contains(string(s), string(ss)))
}
