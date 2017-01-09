package vm

import "testing"

func TestNilEqual(t *testing.T) {
	if !testEqual(NilThunk(), NilThunk()) {
		t.Fail()
	}

	if testEqual(NilThunk(), True) {
		t.Fail()
	}
}
