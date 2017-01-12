package vm

import "testing"

func TestStringEqual(t *testing.T) {
	n := NewString("foo")

	if !testEqual(n, n) {
		t.Fail()
	}

	if testEqual(n, NewString("bar")) {
		t.Fail()
	}
}

func TestStringAdd(t *testing.T) {
	s := "foo"
	st := NewString(s)

	if string(App(Add, st, st).Eval().(stringType)) != s+s {
		t.Fail()
	}
}
