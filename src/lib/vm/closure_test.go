package vm

import "testing"

func TestPartial(t *testing.T) {
	ifFunc := func(ts ...*Thunk) bool {
		b := App(App(Normal(Partial), Normal(If), Normal(False), Normal(True)), ts...)
		return bool(b.EvalStrictly().(Bool))
	}

	if !ifFunc(Normal(True)) {
		t.Fail()
	}
}
