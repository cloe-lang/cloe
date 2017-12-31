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

	th := core.PApp(
		builtinsEnvironment().get("$map"),
		identity,
		core.PApp(core.PApp(builtins.Y, many42)))

	startTimer()

	for i := 0; i < N; i++ {
		if core.NumberType(42) != core.PApp(core.First, th).Eval().(core.NumberType) {
			fail()
		}

		th = core.PApp(core.Rest, th)
	}
}

var identity = core.NewLazyFunction(
	core.NewSignature([]string{"arg"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		return ts[0]
	})

var many42 = core.NewLazyFunction(
	core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		return core.PApp(core.Prepend, core.NewNumber(42), core.PApp(ts[0]))
	})
