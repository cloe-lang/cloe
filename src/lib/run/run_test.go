package run

import (
	"../vm"
	"testing"
)

func TestRunWithNoThunk(t *testing.T) {
	Run([]*vm.Thunk{})
}

func TestRunWithOneThunk(t *testing.T) {
	Run([]*vm.Thunk{vm.PApp(vm.Add, vm.NewNumber(123), vm.NewNumber(456))})
}

func TestRunWithThunks(t *testing.T) {
	th := vm.PApp(vm.Add, vm.NewNumber(123), vm.NewNumber(456))
	Run([]*vm.Thunk{th, th, th, th, th, th, th, th})
}
