package vm

import "testing"

func TestNumberEqual(t *testing.T) {
	n := NewNumber(123)

	if !testEqual(n, n) {
		t.Log(n.Eval(), n.Eval())
		t.Fail()
	}

	if testEqual(n, NewNumber(456)) {
		t.Fail()
	}
}
