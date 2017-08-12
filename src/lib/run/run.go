package run

import (
	"fmt"
	"os"
	"sync"

	"github.com/tisp-lang/tisp/src/lib/compile"
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/systemt"
)

const maxConcurrentOutputs = 256

var sem = make(chan bool, maxConcurrentOutputs)

// Run runs outputs.
func Run(os []compile.Output) {
	go systemt.RunDaemons()

	wg := sync.WaitGroup{}

	for _, v := range os {
		wg.Add(1)

		if v.Expanded() {
			go evalOutputList(v.Value(), &wg, fail)
		} else {
			sem <- true
			go runOutput(v.Value(), &wg, fail)
		}
	}

	wg.Wait()
}

func evalOutputList(t *core.Thunk, parent *sync.WaitGroup, fail func(error)) {
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		parent.Done()
	}()

	for {
		v := core.PApp(core.Equal, t, core.EmptyList).Eval()

		if b, ok := v.(core.BoolType); !ok {
			fail(core.NotBoolError(v).Eval().(core.ErrorType))
		} else if b {
			break
		}

		wg.Add(1)
		sem <- true
		go runOutput(core.PApp(core.First, t), &wg, fail)

		t = core.PApp(core.Rest, t)
	}
}

func runOutput(t *core.Thunk, wg *sync.WaitGroup, fail func(error)) {
	defer func() {
		<-sem
		wg.Done()
	}()

	if err, ok := t.EvalOutput().(core.ErrorType); ok {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
}
