package vm

import "strings"

type StringType string

func NewString(s string) *Thunk {
	return Normal(StringType(s))
}

func (s StringType) equal(e equalable) Object {
	return rawBool(s == e)
}

func (s StringType) merge(ts ...*Thunk) Object {
	return App(NewStrictFunction(func(os ...Object) Object {
		ss := make([]string, len(os))

		for i, o := range os {
			s, ok := o.(StringType)

			if !ok {
				return TypeError(o, "String")
			}

			ss[i] = string(s)
		}

		return StringType(strings.Join(ss, ""))
	}), append([]*Thunk{Normal(s)}, ts...)...)
}

func (s StringType) toList() Object {
	if s == "" {
		return emptyList
	}

	rs := []rune(string(s))

	return cons(
		NewString(string(rs[0])),
		App(ToList, NewString(string(rs[1:]))))
}

// ordered

func (s StringType) less(o ordered) bool {
	return s < o.(StringType)
}
