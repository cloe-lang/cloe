package vm

import "testing"

func TestStringEqual(t *testing.T) {
	n := NewString("foo")

	if !testEqual(n, n) {
		t.Log(n.Eval(), n.Eval())
		t.Fail()
	}

	if testEqual(n, NewString("bar")) {
		t.Fail()
	}
}

func TestStringAdd(t *testing.T) {
	s := "foo"
	ss := NewString(s)

	if string(Add(ss, ss).Eval().(String)) != s+s {
		t.Fail()
	}
}
