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

func NewValueThunk(v types.Object) *Thunk {
	return &Thunk{result: v, state: VALUE}
}

func NewAppThunk(f types.Object, as types.Object) *Thunk {
	return &Thunk{function: f, args: as, state: APP}
}

func (t *Thunk) Wait() {
	t.blackHole.Wait()
}

func (t *Thunk) Lock() bool {
	return t.compareAndSwapState(APP, LOCKED)
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

func (t *Thunk) compareAndSwapState(old, new State) bool {
	return atomic.CompareAndSwapUint32((*uint32)(&t.state), uint32(old), uint32(new))
}

func (t *Thunk) storeState(new State) {
	atomic.StoreUint32((*uint32)(&t.state), uint32(new))
}
