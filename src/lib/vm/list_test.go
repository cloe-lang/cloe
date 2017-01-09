package vm

import "testing"

func TestListEqual(t *testing.T) {
	for _, tss := range [][][]*Thunk{
		{{}, {}},
		{{tTrue}, {tTrue}},
		{{tTrue, tFalse}, {tTrue, tFalse}},
	} {
		if !testEqual(ListThunk(tss[0]...), ListThunk(tss[1]...)) {
			t.Fail()
		}
	}

	for _, tss := range [][][]*Thunk{
		{{}, {tTrue}},
		{{tTrue}, {tFalse}},
		{{tTrue, tTrue}, {tTrue, tTrue, tTrue}},
	} {
		if testEqual(ListThunk(tss[0]...), ListThunk(tss[1]...)) {
			t.Fail()
		}
	}
}
