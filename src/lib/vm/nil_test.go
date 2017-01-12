package vm

import "testing"

func TestNilEqual(t *testing.T) {
	if !testEqual(Nil, Nil) {
		t.Fail()
	}

	if testEqual(Nil, True) {
		t.Fail()
	}
}
