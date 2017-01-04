package vm

import (
	"sync"
	"sync/atomic"
)

type State uint32

const (
	illegal State = iota
	normal
	locked
	abnormal
)

type Thunk struct {
	Result    Object
	function  *Thunk
	args      *Thunk
	state     State
	blackHole sync.WaitGroup
}

func NewValueThunk(o Object) *Thunk {
	return &Thunk{Result: o, state: normal}
}

func NewAppThunk(f *Thunk, args *Thunk) *Thunk {
	t := &Thunk{function: f, args: args, state: abnormal}
	t.blackHole.Add(1)
	return t
}

func (t *Thunk) Eval() { // into WHNF
	if !t.compareAndSwapState(abnormal, locked) {
		// Some goroutine is evaluating this thunk or it has been evaluated already.
		return
	}

	go t.function.Eval()
	go t.args.Eval()

	t.function.Wait()
	t.args.Wait()

	f, ok := t.function.Result.(Callable)

	if !ok {
		panic("Something not callable was called.")
	}

	args, ok := t.args.Result.(Dictionary)

	if !ok {
		panic("Something which is not a dictionary was used as arguments.")
	}

	t.Result = f.Call(args)
	t.function = nil
	t.args = nil

	t.storeState(normal)

	t.blackHole.Done()
}

func (t *Thunk) Wait() {
	t.blackHole.Wait()
}

func (t *Thunk) compareAndSwapState(old, new State) bool {
	return atomic.CompareAndSwapUint32(
		(*uint32)(&t.state),
		uint32(old),
		uint32(new))
}

func (t *Thunk) storeState(new State) {
	atomic.StoreUint32((*uint32)(&t.state), uint32(new))
}
