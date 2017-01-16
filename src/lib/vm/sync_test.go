package vm

import (
	"fmt"
	"testing"
)

func TestSync(t *testing.T) {
	_, ok := App(Sync, App(print, NewNumber(42)), App(print, NewString("OK!"))).Eval().(nilType)

	if !ok {
		t.Fail()
	}
}

var print = NewStrictFunction(func(os ...Object) Object {
	xs := make([]interface{}, len(os))

	for i, o := range os {
		xs[i] = o
	}

	fmt.Println(xs...)
	return Nil
})
