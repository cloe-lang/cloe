package run

import (
	"fmt"
	"github.com/raviqqe/tisp/src/lib/compile"
	"github.com/raviqqe/tisp/src/lib/core"
	"os"
	"sync"
)

const maxConcurrentOutputs = 256

var sem = make(chan bool, maxConcurrentOutputs)

// Run runs outputs.
func Run(os []compile.Output) {
	// TODO: Ensure results are Outputs.

	wg := sync.WaitGroup{}

	for _, o := range os {
		wg.Add(1)

		if o.Expanded() {
			go func() {
				evalOutputList(o.Value())
				wg.Done()
			}()
		} else {
			sem <- true
			go runOutput(o.Value(), &wg)
		}
	}

	wg.Wait()
}

func evalOutputList(t *core.Thunk) {
	wg := sync.WaitGroup{}

	for {
		o := core.PApp(core.Equal, t, core.EmptyList).Eval()

		if b, ok := o.(core.BoolType); !ok {
			failOnError(core.NotBoolError(o).Eval().(core.ErrorType))
		} else if b {
			break
		}

		wg.Add(1)
		sem <- true
		go runOutput(core.PApp(core.First, t), &wg)

		t = core.PApp(core.Rest, t)
	}

	wg.Wait()
}

func runOutput(t *core.Thunk, wg *sync.WaitGroup) {
	if err, ok := t.Eval().(core.ErrorType); ok {
		failOnError(err)
	}

	<-sem
	wg.Done()
}

func failOnError(err core.ErrorType) {
	fmt.Fprint(os.Stderr, err.Lines())
	os.Exit(1)
}
