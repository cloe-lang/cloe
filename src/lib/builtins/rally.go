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
	func(vs ...core.Value) core.Value {
		c := make(chan core.Value, valueChannelCapacity)

		systemt.Daemonize(func() {
			l, err := core.EvalList(vs[0])

			if err != nil {
				c <- err
				return
			}

			sem := make(chan bool, maxConcurrency)

			for !l.Empty() {
				sem <- true
				go func(v core.Value) {
					c <- core.EvalPure(v)
					<-sem
				}(l.First())

				l, err = core.EvalList(l.Rest())

				if err != nil {
					c <- err
					break
				}
			}

			// HACK: Wait for other goroutines to put elements in a value channel
			// for a while. This is only for unit test.
			time.Sleep(channelCloseDuration)
			c <- nil
		})

		return core.PApp(core.PApp(Y, core.NewLazyFunction(
			core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
			func(vs ...core.Value) core.Value {
				v := <-c

				if v == nil {
					return core.EmptyList
				} else if err, ok := v.(core.ErrorType); ok {
					return err
				}

				return core.StrictPrepend([]core.Value{v}, core.PApp(vs[0]))
			})))
	})
