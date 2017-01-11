package vm

import "testing"

func TestBoolEqual(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{True, True},
		{False, False},
	} {
		if !testEqual(ts...) {
			t.Fail()
		}
	}

	for _, ts := range [][]*Thunk{
		{True, False},
		{False, True},
	} {
		if testEqual(ts...) {
			t.Fail()
		}
	}
}
