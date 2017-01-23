package ir

import "testing"

func TestRun(t *testing.T) {
	Run([]interface{}{
		NewLetFunction("foo", []interface{}{"+", 0, 1}),
		NewOutput([]interface{}{"write", []interface{}{"foo", "123", "456"}}),
	})
}
