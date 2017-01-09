package vm

import "testing"

func TestPartial(t *testing.T) {
	ifFunc := func(ts ...*Thunk) bool {
		f := Partial(Normal(If), Normal(False), Normal(True)).(Callable)
		return bool(f.Call(ts...).(Bool))
	}

	if !ifFunc(Normal(True)) {
		t.Fail()
	}
}
