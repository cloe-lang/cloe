package vm

import "testing"

func TestNumberEqual(t *testing.T) {
	n := NewNumber(123)

	if !testEqual(n, n) {
		t.Fail()
	}

	if testEqual(n, NewNumber(456)) {
		t.Fail()
	}
}

func TestNumberAdd(t *testing.T) {
	var n1, n2 float64 = 123, 456

	if float64(Add(NewNumber(n1), NewNumber(n2)).Eval().(Number)) != n1+n2 {
		t.Fail()
	}
}
