package vm

import (
	"github.com/mediocregopher/seq"
	"hash/crc32"
	"strings"
)

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
			return typeError(o, "String")
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

// seq.Setable

func (s stringType) Hash(i uint32) uint32 {
	// TODO: Need to add i?
	return crc32.ChecksumIEEE([]byte(s)) % seq.ARITY
}

func (s1 stringType) Equal(o interface{}) bool {
	s2, ok := o.(stringType)

	if !ok {
		return false
	}

	return s1 == s2
}
