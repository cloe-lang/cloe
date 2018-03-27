package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDump(t *testing.T) {
	for _, c := range []struct {
		argument Value
		answer   StringType
	}{
		{NewString("foo"), `"foo"`},
		{NewList(NewString("foo")), `["foo"]`},
		{NewDictionary([]KeyValue{{NewString("foo"), NewString("bar")}}), `{"foo" "bar"}`},
	} {
		assert.Equal(t, c.answer, EvalPure(PApp(Dump, c.argument)).(StringType))
	}
}

func TestDumpError(t *testing.T) {
	for _, v := range []Value{
		DummyError,
		NewList(DummyError),
	} {
		v := EvalPure(PApp(Dump, v))
		t.Log(v)

		_, ok := v.(*ErrorType)
		assert.True(t, ok)
	}
}

func TestStrictDumpPanic(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	StrictDump(nil)
}

func TestInternalStrictDumpFail(t *testing.T) {
	for _, v := range []Value{
		NewList(DummyError),
		NewDictionary([]KeyValue{{Nil, DummyError}}),
	} {
		_, err := StrictDump(EvalPure(v))
		assert.NotNil(t, err)
	}
}

func TestStrictDump(t *testing.T) {
	for _, v := range []Value{
		Nil,
		True,
		False,
		EmptyList,
		EmptyDictionary,
		NewNumber(42),
		NewString("foo"),
	} {
		s, err := StrictDump(EvalPure(v))
		assert.NotEqual(t, "", s)
		assert.Nil(t, err)
	}
}

func TestIdentity(t *testing.T) {
	for _, v := range []Value{Nil, NewNumber(42), True, False, NewString("foo")} {
		assert.True(t, testEqual(PApp(identity, v), v))
	}
}

func TestTypeOf(t *testing.T) {
	for _, test := range []struct {
		typ   string
		thunk Value
	}{
		{"nil", Nil},
		{"list", EmptyList},
		{"list", EmptyList},
		{"bool", True},
		{"number", NewNumber(123)},
		{"string", NewString("foo")},
		{"dict", EmptyDictionary},
		{"function", identity},
		{"function", PApp(Partial, identity)},
	} {
		v := EvalPure(PApp(TypeOf, test.thunk))
		t.Log(v)
		assert.Equal(t, test.typ, string(v.(StringType)))
	}
}

func TestTypeOfError(t *testing.T) {
	for _, v := range []Value{
		NewError("MyError", "This is error."),
		PApp(impureFunction, NewNumber(42)),
	} {
		v = EvalPure(PApp(TypeOf, v))
		_, ok := v.(*ErrorType)
		t.Logf("%#v\n", v)
		assert.True(t, ok)
	}
}
