package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	for _, test := range []struct {
		typ   string
		thunk *Thunk
	}{
		{"nil", Nil},
		{"list", NewList()},
		{"list", EmptyList},
		{"bool", True},
		{"number", NewNumber(123)},
		{"string", NewString("foo")},
		{"dict", NewDictionary(nil, nil)},
		{"error", NewError("MyError", "This is error.")},
	} {
		v := PApp(TypeOf, test.thunk).Eval()
		t.Log(v)
		assert.Equal(t, test.typ, string(v.(StringType)))
	}
}

func TestTypeOfFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		} else if _, ok := r.(error); !ok {
			t.Fail()
		}
	}()

	PApp(TypeOf, Normal("foo")).Eval()
}
