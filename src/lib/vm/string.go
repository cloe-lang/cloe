package vm

import "strings"

type stringType string

func NewString(s string) *Thunk {
	return Normal(stringType(s))
}

func (s stringType) equal(e equalable) Object {
	return rawBool(s == e)
}

var Concat = NewStrictFunction(func(os ...Object) Object {
	ss := make([]string, len(os))

	for i, o := range os {
		s, ok := o.(stringType)

		if !ok {
			return TypeError(o, "String")
		}

		ss[i] = string(s)
	}

	return stringType(strings.Join(ss[:], ""))
})

func (s stringType) toList() Object {
	if s == "" {
		return emptyList
	}

	rs := []rune(string(s))

	return cons(
		NewString(string(rs[0])),
		App(ToList, NewString(string(rs[1:]))))
}

// ordered

func (s stringType) less(o ordered) bool {
	return s < o.(stringType)
}
