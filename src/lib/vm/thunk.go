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
	app
)

type Thunk struct {
	Result    Object
	function  *Thunk
	args      *Thunk
	state     State
	blackHole sync.WaitGroup
}

func Normal(o Object) *Thunk {
	return &Thunk{Result: o, state: normal}
}

func NormalApp(f Function, args List) *Thunk {
	return App(Normal(f), Normal(args))
}

func App(f *Thunk, args *Thunk) *Thunk {
	t := &Thunk{function: f, args: args, state: app}
	t.blackHole.Add(1)
	return t
}

func (t *Thunk) Eval() Object { // into WHNF
	if t.compareAndSwapState(app, locked) {
		go t.function.Eval()
		go t.args.Eval()

		f, ok := t.function.Eval().(Callable)

		if !ok {
			panic("Something not callable was called.")
		}

		args, ok := t.args.Eval().(List)

		if !ok {
			panic("Something which is not a list was used as arguments.")
		}

		t.Result = f.Call(args).Eval()

		t.function = nil
		t.args = nil

		t.storeState(normal)

		t.blackHole.Done()
	} else {
		t.blackHole.Wait()
	}

	return t.Result
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
