package vm

import (
	"github.com/mediocregopher/seq"
	"hash/crc32"
)

type stringType string

func NewString(s string) *Thunk {
	return Normal(stringType(s))
}

func (s stringType) equal(e equalable) Object {
	return rawBool(s == e)
}

func (s stringType) add(a addable) addable {
	return s + a.(stringType)
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
