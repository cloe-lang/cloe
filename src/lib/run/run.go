package run

import (
	"fmt"
	"os"
	"sync"

	"github.com/coel-lang/coel/src/lib/compile"
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
)

const maxConcurrentEffects = 256

var sem = make(chan bool, maxConcurrentEffects)

// Run runs effects.
func Run(os []compile.Effect) {
	go systemt.RunDaemons()

	wg := sync.WaitGroup{}

	for _, v := range os {
		wg.Add(1)

		if v.Expanded() {
			go evalEffectList(v.Value(), &wg, fail)
		} else {
			sem <- true
			go runEffect(v.Value(), &wg, fail)
		}
	}

	wg.Wait()
}

func evalEffectList(v core.Value, parent *sync.WaitGroup, fail func(error)) {
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		parent.Done()
	}()

	for {
		if b, err := core.EvalBool(core.EvalPure(core.PApp(core.Equal, v, core.EmptyList))); err != nil {
			fail(err.(core.ErrorType))
		} else if *b {
			break
		}

		wg.Add(1)
		sem <- true
		go runEffect(core.PApp(core.First, v), &wg, fail)

		v = core.PApp(core.Rest, v)
	}
}

func runEffect(v core.Value, wg *sync.WaitGroup, fail func(error)) {
	defer func() {
		<-sem
		wg.Done()
	}()

	if err, ok := core.EvalImpure(v).(core.ErrorType); ok {
		fail(err)
	}
}

var fail = failWithExit(os.Exit)

func failWithExit(exit func(int)) func(error) {
	return func(err error) {
		_, err = fmt.Fprint(os.Stderr, err)

		if err != nil {
			panic(err)
		}

		exit(1)
	}
}
