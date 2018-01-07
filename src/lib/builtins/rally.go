package builtins

import (
	"time"

	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
)

const maxConcurrency = 256
const valueChannelCapacity = 1024
const channelCloseDuration = 100 * time.Millisecond

// Rally sorts arguments by time.
var Rally = core.NewLazyFunction(
	core.NewSignature(nil, nil, "args", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		vs := make(chan core.Value, valueChannelCapacity)

		systemt.Daemonize(func() {
			l, err := ts[0].EvalList()

			if err != nil {
				vs <- err
				return
			}

			sem := make(chan bool, maxConcurrency)

			for !l.Empty() {
				sem <- true
				go func(t *core.Thunk) {
					vs <- t.Eval()
					<-sem
				}(l.First())

				l, err = l.Rest().EvalList()

				if err != nil {
					vs <- err
					break
				}
			}

			// HACK: Wait for other goroutines to put elements in a value channel
			// for a while. This is only for unit test.
			time.Sleep(channelCloseDuration)
			vs <- nil
		})

		return core.PApp(core.PApp(Y, core.NewLazyFunction(
			core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
			func(ts ...*core.Thunk) core.Value {
				v := <-vs

				if v == nil {
					return core.EmptyList
				} else if err, ok := v.(core.ErrorType); ok {
					return err
				}

				return core.PApp(core.Prepend, core.Normal(v), core.PApp(ts[0]))
			})))
	})
