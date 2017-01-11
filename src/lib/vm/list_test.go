package vm

import "testing"

func TestListEqual(t *testing.T) {
	for _, tss := range [][][]*Thunk{
		{{}, {}},
		{{True}, {True}},
		{{True, False}, {True, False}},
	} {
		if !testEqual(ListThunk(tss[0]...), ListThunk(tss[1]...)) {
			t.Fail()
		}
	}

	for _, tss := range [][][]*Thunk{
		{{}, {True}},
		{{True}, {False}},
		{{True, True}, {True, True, True}},
	} {
		if testEqual(ListThunk(tss[0]...), ListThunk(tss[1]...)) {
			t.Fail()
		}
	}
}
