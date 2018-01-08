package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDump(t *testing.T) {
	for _, c := range []struct {
		argument *Thunk
		answer   StringType
	}{
		{NewString("foo"), `"foo"`},
		{NewList(NewString("foo")), `["foo"]`},
		{NewDictionary([]KeyValue{{NewString("foo"), NewString("bar")}}), `{"foo" "bar"}`},
	} {
		assert.Equal(t, c.answer, PApp(Dump, c.argument).Eval().(StringType))
	}
}

func TestDumpError(t *testing.T) {
	for _, th := range []*Thunk{
		OutOfRangeError(),
		NewList(OutOfRangeError()),
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

	StrictDump(nil)
}

func TestInternalStrictDumpPanicWithEffect(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	StrictDump(effectType{Nil})
}

func TestInternalStrictDumpFail(t *testing.T) {
	for _, th := range []*Thunk{
		NewList(OutOfRangeError()),
		NewDictionary([]KeyValue{{Nil, OutOfRangeError()}}),
	} {
		_, err := StrictDump(th.Eval())
		assert.NotNil(t, err)
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
	} {
		s, err := StrictDump(th.Eval())
		assert.NotEqual(t, "", s)
		assert.Nil(t, err)
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
		{"function", identity},
		{"function", PApp(Partial, identity)},
	} {
		v := PApp(TypeOf, test.thunk).Eval()
		t.Log(v)
		assert.Equal(t, test.typ, string(v.(StringType)))
	}
}

func TestTypeOfError(t *testing.T) {
	for _, th := range []*Thunk{
		NewError("MyError", "This is error."),
		PApp(impureFunction, NewNumber(42)),
	} {
		v := PApp(TypeOf, th).Eval()
		_, ok := v.(ErrorType)
		t.Log(v)
		assert.True(t, ok)
	}
}

func TestTypeOfFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	PApp(TypeOf, Normal("foo")).Eval()
}
