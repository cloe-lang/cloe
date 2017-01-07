package vm

import "testing"

func TestNilEqual(t *testing.T) {
	if !testEqual(NewNil(), NewNil()) {
		t.Fail()
	}

	if testEqual(NewNil(), True) {
		t.Fail()
	}
}
