package run

import (
	"../core"
	"testing"
)

func TestRunWithNoThunk(t *testing.T) {
	Run([]*core.Thunk{})
}

func TestRunWithOneThunk(t *testing.T) {
	Run([]*core.Thunk{core.PApp(core.Add, core.NewNumber(123), core.NewNumber(456))})
}

func TestRunWithThunks(t *testing.T) {
	th := core.PApp(core.Add, core.NewNumber(123), core.NewNumber(456))
	Run([]*core.Thunk{th, th, th, th, th, th, th, th})
}
