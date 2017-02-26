package core

import "strings"

type StringType string

func NewString(s string) *Thunk {
	return Normal(StringType(s))
}

func (s StringType) equal(e equalable) Object {
	return rawBool(s == e)
}

func (s StringType) merge(ts ...*Thunk) Object {
	ts = append([]*Thunk{Normal(s)}, ts...)

	for _, t := range ts {
		go t.Eval()
	}

	ss := make([]string, len(ts))

	for i, t := range ts {
		o := t.Eval()
		s, ok := o.(StringType)

		if !ok {
			return TypeError(o, "String")
		}

		ss[i] = string(s)
	}

	return StringType(strings.Join(ss, ""))
}

func (s StringType) toList() Object {
	if s == "" {
		return emptyList
	}

	rs := []rune(string(s))

	return cons(
		NewString(string(rs[0])),
		PApp(ToList, NewString(string(rs[1:]))))
}

// ordered

func (s StringType) less(o ordered) bool {
	return s < o.(StringType)
}
