package vm

import "testing"

var n1, n2 float64 = 123, 456

func TestNumberEqual(t *testing.T) {
	n := NumberThunk(123)

	if !testEqual(n, n) {
		t.Fail()
	}

	if testEqual(n, NumberThunk(456)) {
		t.Fail()
	}
}

func TestNumberAdd(t *testing.T) {
	if float64(App(Add, NumberThunk(n1), NumberThunk(n2)).Eval().(Number)) != n1+n2 {
		t.Fail()
	}
}

func TestNumberSub(t *testing.T) {
	if float64(App(Sub, NumberThunk(n1), NumberThunk(n2)).Eval().(Number)) != n1-n2 {
		t.Fail()
	}
}

func TestNumberMult(t *testing.T) {
	if float64(App(Mult, NumberThunk(n1), NumberThunk(n2)).Eval().(Number)) != n1*n2 {
		t.Fail()
	}
}

func TestNumberDiv(t *testing.T) {
	if float64(App(Div, NumberThunk(n1), NumberThunk(n2)).Eval().(Number)) != n1/n2 {
		t.Fail()
	}
}
