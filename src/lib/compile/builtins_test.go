package compile

import (
	"testing"
	"time"

	"github.com/coel-lang/coel/src/lib/builtins"
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinsEnvironment(t *testing.T) {
	builtinsEnvironment()
}

func TestCompileBuiltinModule(t *testing.T) {
	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) x)`)
}

func TestCompileBuiltinModuleWithInvalidSyntax(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) x`)
}

func TestCompileBuiltinModuleWithInvalidSource(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) y)`)
}

func TestReduce(t *testing.T) {
	f := builtinsEnvironment().get("$reduce")

	for _, ts := range [][2]*core.Thunk{
		{
			core.PApp(f, core.Add, core.NewList(core.NewNumber(1), core.NewNumber(2), core.NewNumber(3))),
			core.NewNumber(6),
		},
		{
			core.PApp(f, core.Sub, core.NewList(core.NewNumber(1), core.NewNumber(2), core.NewNumber(3))),
			core.NewNumber(-4),
		},
	} {
		t.Log(core.PApp(core.ToString, ts[0]).Eval())
		assert.True(t, bool(core.PApp(core.Equal, ts[0], ts[1]).Eval().(core.BoolType)))
	}
}

func TestReduceError(t *testing.T) {
	f := builtinsEnvironment().get("$reduce")

	for _, th := range []*core.Thunk{
		core.PApp(f, core.Add, core.EmptyList),
		core.PApp(f, core.IsOrdered, core.EmptyList),
	} {
		_, ok := th.Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestMapOrder(t *testing.T) {
	b := func(N int) float64 {
		var start time.Time
		benchmarkMap(N, func() { start = time.Now() }, t.Fail)
		return time.Since(start).Seconds()
	}

	r := b(10000) / b(2000)
	t.Log(r)
	assert.True(t, 4 < r && r < 6)
}

func BenchmarkMap(b *testing.B) {
	benchmarkMap(b.N, b.ResetTimer, b.Fail)
}

func benchmarkMap(N int, startTimer, fail func()) {
	go systemt.RunDaemons()

	th := core.PApp(builtinsEnvironment().get("$map"), identity, many42())

	startTimer()

	for i := 0; i < N; i++ {
		if core.NumberType(42) != core.PApp(core.First, th).Eval().(core.NumberType) {
			fail()
		}

		th = core.PApp(core.Rest, th)
	}
}

func BenchmarkGoMap(b *testing.B) {
	th := many42()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if core.NumberType(42) != core.PApp(core.First, th).Eval().(core.NumberType) {
			b.Fail()
		}

		th = core.PApp(core.Rest, th)
	}
}

var identity = core.NewLazyFunction(
	core.NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		return ts[0]
	})

func many42() *core.Thunk {
	return core.PApp(core.PApp(builtins.Y, core.NewLazyFunction(
		core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
		func(ts ...*core.Thunk) core.Value {
			return core.PApp(core.Prepend, core.NewNumber(42), core.PApp(ts[0]))
		})))
}
