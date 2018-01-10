package compile

import (
	"math"
	"math/rand"
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
	for _, f := range []func(){
		func() { compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) y)`) },
		func() { compileBuiltinModule(newEnvironment(testFallback), "", `(import "foo")`) },
	} {
		assert.Panics(t, f)
	}
}

func TestReduce(t *testing.T) {
	f := builtinsEnvironment().get("$reduce")

	for _, ts := range [][2]core.Value{
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

	for _, th := range []core.Value{
		core.PApp(f, core.Add, core.EmptyList),
		core.PApp(f, core.IsOrdered, core.EmptyList),
	} {
		_, ok := th.Eval().(core.ErrorType)
		assert.True(t, ok)
	}
}

func TestFilter(t *testing.T) {
	f := builtinsEnvironment().get("$filter")

	for _, ts := range [][2]core.Value{
		{
			core.PApp(f, core.IsOrdered, core.EmptyList),
			core.EmptyList,
		},
		{
			core.PApp(f, core.IsOrdered, core.NewList(core.NewString("foo"))),
			core.NewList(core.NewString("foo")),
		},
		{
			core.PApp(f,
				core.IsOrdered,
				core.NewList(core.NewNumber(42), core.EmptyDictionary, core.Nil, core.EmptyList)),
			core.NewList(core.NewNumber(42), core.EmptyList),
		},
	} {
		t.Log(core.PApp(core.ToString, ts[0]).Eval())
		assert.True(t, bool(core.PApp(core.Equal, ts[0], ts[1]).Eval().(core.BoolType)))
	}
}

func BenchmarkFilter(b *testing.B) {
	go systemt.RunDaemons()

	f := builtinsEnvironment().get("$filter")
	l := randomNumberList(10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		core.PApp(core.ToString, core.PApp(f, core.IsOrdered, l)).Eval()
	}
}

func BenchmarkFilterBare(b *testing.B) {
	go systemt.RunDaemons()

	f := core.PApp(builtins.Y, core.NewLazyFunction(
		core.NewSignature([]string{"me", "func", "list"}, nil, "", nil, nil, ""),
		func(ts ...core.Value) core.Value {
			f := ts[1]
			l := ts[2]

			return core.PApp(core.If,
				core.PApp(core.Equal, l, core.EmptyList),
				core.EmptyList,
				core.PApp(core.Prepend,
					core.PApp(f, core.PApp(core.First, l)),
					core.PApp(ts[0], f, core.PApp(core.Rest, l))))
		}))

	l := randomNumberList(10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		core.PApp(core.ToString, core.PApp(f, core.IsOrdered, l)).Eval()
	}
}

func TestSort(t *testing.T) {
	go systemt.RunDaemons()

	for _, ts := range [][2]core.Value{
		{
			core.EmptyList,
			core.EmptyList,
		},
		{
			core.NewList(core.NewNumber(42)),
			core.NewList(core.NewNumber(42)),
		},
		{
			core.NewList(core.NewNumber(2), core.NewNumber(1)),
			core.NewList(core.NewNumber(1), core.NewNumber(2)),
		},
		{
			core.NewList(core.NewNumber(1), core.NewNumber(1)),
			core.NewList(core.NewNumber(1), core.NewNumber(1)),
		},
		{
			core.NewList(core.NewNumber(3), core.NewNumber(2), core.NewNumber(1)),
			core.NewList(core.NewNumber(1), core.NewNumber(2), core.NewNumber(3)),
		},
		{
			core.NewList(core.NewNumber(2), core.NewNumber(3), core.NewNumber(1), core.NewNumber(-123)),
			core.NewList(core.NewNumber(-123), core.NewNumber(1), core.NewNumber(2), core.NewNumber(3)),
		},
	} {
		th := core.PApp(builtinsEnvironment().get("sort"), ts[0])
		t.Log(core.PApp(core.ToString, th).Eval())
		assert.True(t, bool(core.PApp(core.Equal, th, ts[1]).Eval().(core.BoolType)))
	}
}

func TestSortError(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.App(
		builtinsEnvironment().get("$sort"),
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewList(core.NewNumber(42)), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("less", builtins.LessEq),
			},
			nil)).Eval().(core.ErrorType)

	assert.True(t, ok)
}

func TestSortWithBigLists(t *testing.T) {
	for i := 0; i < 4; i++ {
		benchmarkSort(int(math.Pow10(i)), 1, func() {})
	}
}

func BenchmarkSort100(b *testing.B) {
	benchmarkSort(100, b.N, b.ResetTimer)
}

func BenchmarkSort1000(b *testing.B) {
	benchmarkSort(1000, b.N, b.ResetTimer)
}

func BenchmarkSort10000(b *testing.B) {
	benchmarkSort(10000, b.N, b.ResetTimer)
}

func benchmarkSort(size, N int, resetTimer func()) {
	go systemt.RunDaemons()

	f := builtinsEnvironment().get("$sort")
	l := randomNumberList(size)

	resetTimer()

	for i := 0; i < N; i++ {
		core.PApp(core.First, core.PApp(f, l)).Eval()
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
	func(ts ...core.Value) core.Value {
		return ts[0]
	})

func many42() core.Value {
	return core.PApp(core.PApp(builtins.Y, core.NewLazyFunction(
		core.NewSignature([]string{"me"}, nil, "", nil, nil, ""),
		func(ts ...core.Value) core.Value {
			return core.PApp(core.Prepend, core.NewNumber(42), core.PApp(ts[0]))
		})))
}

func randomNumberList(n int) core.Value {
	r := rand.New(rand.NewSource(42))
	ts := make([]core.Value, n)

	for i := range ts {
		ts[i] = core.NewNumber(r.Float64())
	}

	l := core.NewList(ts...)
	l.Eval()
	return l
}
