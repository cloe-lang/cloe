package vm

import "testing"

func TestStringEqual(t *testing.T) {
	n := StringThunk("foo")

	if !testEqual(n, n) {
		t.Fail()
	}

	if testEqual(n, StringThunk("bar")) {
		t.Fail()
	}
}

func TestStringAdd(t *testing.T) {
	s := "foo"
	ss := StringThunk(s)

	if string(Add(ss, ss).Eval().(String)) != s+s {
		t.Fail()
	}
}
