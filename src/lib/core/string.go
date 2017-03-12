package core

import "strings"

// StringType represents strings in the language.
type StringType string

// NewString creates a string in the language from one in Go.
func NewString(s string) *Thunk {
	return Normal(StringType(s))
}

func (s StringType) equal(e equalable) Value {
	return rawBool(s == e)
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
