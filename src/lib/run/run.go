package run

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/core"
	"sync"
)

const maxConcurrentOutputs = 256

var outSem = make(chan bool, maxConcurrentOutputs)

func Run(os []compile.Output) {
	// TODO: Ensure results are Outputs.

	wg := sync.WaitGroup{}

	for _, o := range os {
		if o.Expanded() {
			wg.Add(1)
			go func() {
				evalOutputList(o.Value())
				wg.Done()
			}()

			continue
		}

		wg.Add(1)
		outSem <- true
		go func() {
			o.Value().Eval()
			<-outSem
			wg.Done()
		}()
	}

	wg.Wait()
}

func evalOutputList(t *core.Thunk) {
	wg := sync.WaitGroup{}

	for {
		o := core.PApp(core.Equal, t, core.EmptyList).Eval()

		if b, ok := o.(core.BoolType); !ok {
			// TODO: (write error)
			fmt.Println(o)
			break
		} else if b {
			break
		}

		wg.Add(1)
		outSem <- true
		go func(t *core.Thunk) {
			core.PApp(core.First, t).Eval() // TODO: Check error
			<-outSem
			wg.Done()
		}(t)

		t = core.PApp(core.Rest, t)
	}

	wg.Wait()
}
