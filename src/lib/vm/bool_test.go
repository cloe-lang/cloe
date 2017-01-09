package vm

import "testing"

var tTrue, tFalse = Normal(True), Normal(False)

func TestBoolEqual(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{tTrue, tTrue},
		{tFalse, tFalse},
	} {
		if !testEqual(ts...) {
			t.Fail()
		}
	}

	for _, ts := range [][]*Thunk{
		{tTrue, tFalse},
		{tFalse, tTrue},
	} {
		if testEqual(ts...) {
			t.Fail()
		}
	}
}
