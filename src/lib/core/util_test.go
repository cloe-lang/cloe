package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDumpFail(t *testing.T) {
	for _, th := range []*Thunk{
		OutOfRangeError(),
	} {
		v := PApp(Dump, th).Eval()
		t.Log(v)

		_, ok := v.(ErrorType)
		assert.True(t, ok)
	}
}

func TestInternalStrictDumpPanic(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	strictDump(nil)
}

func TestInternalStrictDumpFail(t *testing.T) {
	for _, th := range []*Thunk{
		NewList(OutOfRangeError()),
		NewDictionary([]KeyValue{{Nil, OutOfRangeError()}})} {
		_, err := strictDump(th.Eval())
		assert.NotEqual(t, (*Thunk)(nil), err)
	}
}

func TestStrictDump(t *testing.T) {
	for _, th := range []*Thunk{
		Nil,
		True,
		False,
		EmptyList,
		EmptyDictionary,
		NewNumber(42),
		NewString("foo"),
		newEffect(Nil),
	} {
		s, err := StrictDump(th.Eval())
		assert.NotEqual(t, "", s)
		assert.Equal(t, (*Thunk)(nil), err)
	}
}

func TestEqualFail(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{OutOfRangeError(), Nil},
		{True, True, OutOfRangeError()},
	} {
		_, ok := PApp(Equal, ts...).Eval().(ErrorType)
		assert.True(t, ok)
	}
}

func TestIdentity(t *testing.T) {
	for _, th := range []*Thunk{Nil, NewNumber(42), True, False, NewString("foo")} {
		assert.True(t, testEqual(PApp(identity, th), th))
	}
}

func TestTypeOf(t *testing.T) {
	for _, test := range []struct {
		typ   string
		thunk *Thunk
	}{
		{"nil", Nil},
		{"list", EmptyList},
		{"list", EmptyList},
		{"bool", True},
		{"number", NewNumber(123)},
		{"string", NewString("foo")},
		{"dict", EmptyDictionary},
		{"error", NewError("MyError", "This is error.")},
		{"function", identity},
		{"function", PApp(Partial, identity)},
	} {
		v := PApp(TypeOf, test.thunk).Eval()
		t.Log(v)
		assert.Equal(t, test.typ, string(v.(StringType)))
	}
}

func TestTypeOfError(t *testing.T) {
	v := PApp(impureFunction, NewNumber(42)).Eval()
	_, ok := v.(ErrorType)
	t.Log(v)
	assert.True(t, ok)
}

func TestTypeOfFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	PApp(TypeOf, Normal("foo")).Eval()
}
