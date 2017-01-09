package vm

import "testing"

var tTrue, tFalse = Normal(True), Normal(False)

func TestBoolEqual(t *testing.T) {
	for _, ts := range [][]*Thunk{
		{tTrue},
		{tFalse},
		{tTrue, tTrue},
		{tFalse, tFalse},
		{tTrue, tTrue, tTrue},
	} {
		if !testEqual(ts...) {
			t.Fail()
		}
	}

	for _, ts := range [][]*Thunk{
		{tTrue, tFalse},
		{tFalse, tTrue},
		{tTrue, tTrue, tFalse},
	} {
		if testEqual(ts...) {
			t.Fail()
		}
	}
}
