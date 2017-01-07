package vm

import "testing"

func TestPartial(t *testing.T) {
	and := func(ts ...*Thunk) bool {
		return bool(App(Partial(Normal(And), True, True), ts...).Eval().(Bool))
	}

	if !and(True) {
		t.Fail()
	}
}
