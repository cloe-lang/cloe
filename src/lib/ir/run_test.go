package ir

import (
	"../vm"
	"testing"
)

func TestRun(t *testing.T) {
	Run([]interface{}{
		NewLetFunction("foo", vm.NewSimpleSignature("x", "y"), []interface{}{"+", 0, 1}),
		NewOutput([]interface{}{"write", []interface{}{"foo", "123", "456"}}),
	})
}
