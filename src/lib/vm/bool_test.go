package vm

import "testing"

func TestNewBool(t *testing.T) {
	for _, b := range []bool{true, false} {
		if bool(NewBool(b).Eval().(Bool)) != b {
			t.Fail()
		}
	}
}

func TestAnd(t *testing.T) {
	and := func(ts ...*Thunk) bool {
		return bool(And(ts...).Eval().(Bool))
	}

	if !and(True, True) {
		t.Fail()
	}

	for _, ts := range [][]*Thunk{{False, False}, {True, False}, {False, True}} {
		if and(ts...) {
			t.Fail()
		}
	}
}
