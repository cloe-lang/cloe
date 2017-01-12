package vm

import "testing"

var n1, n2 float64 = 123, 456

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
	if float64(App(Add, NewNumber(n1), NewNumber(n2)).Eval().(numberType)) != n1+n2 {
		t.Fail()
	}
}

func TestNumberSub(t *testing.T) {
	if float64(App(Sub, NewNumber(n1), NewNumber(n2)).Eval().(numberType)) != n1-n2 {
		t.Fail()
	}
}

func TestNumberMult(t *testing.T) {
	if float64(App(Mult, NewNumber(n1), NewNumber(n2)).Eval().(numberType)) != n1*n2 {
		t.Fail()
	}
}

func TestNumberDiv(t *testing.T) {
	if float64(App(Div, NewNumber(n1), NewNumber(n2)).Eval().(numberType)) != n1/n2 {
		t.Fail()
	}
}
