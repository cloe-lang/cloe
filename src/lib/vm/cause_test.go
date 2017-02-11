package vm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCause(t *testing.T) {
	_, ok := App(Cause, App(print, NewNumber(42)), App(print, NewString("OK!"))).Eval().(NilType)
	assert.True(t, ok)
}

var print = NewStrictFunction(func(os ...Object) Object {
	xs := make([]interface{}, len(os))

	for i, o := range os {
		xs[i] = o
	}

	fmt.Println(xs...)
	return Nil
})
