package vm

import (
	"../types"
	"sync"
	"sync/atomic"
)

type State uint32

const (
	VALUE State = iota
	LOCKED
	APP
)

type Thunk struct {
	result, function, args types.Object
	state                  State
	blackHole              sync.WaitGroup
}

func (t *Thunk) Wait() {
	t.blackHole.Wait()
}

func (t *Thunk) SaveResult(o types.Object) {
	if t.state != LOCKED {
		panic("Thunk is not locked yet.")
	}

	t.result = o
	t.function = nil
	t.args = nil

	t.storeState(VALUE)
}

func (t *Thunk) storeState(new State) {
	atomic.StoreUint32((*uint32)(&t.state), uint32(new))
}
