package vm

import "testing"

func TestPartial(t *testing.T) {
	ifFunc := func(ts ...*Thunk) bool {
		b := App(App(Normal(Partial), Normal(If), False, True), ts...)
		return bool(b.Eval().(boolType))
	}

	if !ifFunc(True) {
		t.Fail()
	}
}
