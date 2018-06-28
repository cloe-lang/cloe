package compile

import (
	"math"
	"math/rand"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/builtins"
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/systemt"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinsEnvironment(t *testing.T) {
	builtinsEnvironment()
}

func TestCompileBuiltinModule(t *testing.T) {
	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) x)`)
}

func TestCompileBuiltinModuleWithInvalidSyntax(t *testing.T) {
	assert.Panics(t, func() {
		compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) x`)
	})
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

	for _, vs := range [][2]core.Value{
		{
			core.PApp(f, core.Add, core.NewList(core.NewNumber(1), core.NewNumber(2), core.NewNumber(3))),
			core.NewNumber(6),
		},
		{
			core.PApp(f, core.Sub, core.NewList(core.NewNumber(1), core.NewNumber(2), core.NewNumber(3))),
			core.NewNumber(-4),
		},
	} {
		t.Log(core.EvalPure(core.PApp(core.ToString, vs[0])))
		assert.True(t, bool(*core.EvalPure(core.PApp(core.Equal, vs[0], vs[1])).(*core.BooleanType)))
	}
}

func TestReduceError(t *testing.T) {
	f := builtinsEnvironment().get("$reduce")

	for _, v := range []core.Value{
		core.PApp(f, core.Add, core.EmptyList),
		core.PApp(f, core.IsOrdered, core.EmptyList),
	} {
		_, ok := core.EvalPure(v).(*core.ErrorType)
		assert.True(t, ok)
	}
}

func TestFilter(t *testing.T) {
	f := builtinsEnvironment().get("$filter")

	for _, vs := range [][2]core.Value{
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
		t.Log(core.EvalPure(core.PApp(core.ToString, vs[0])))
		assert.True(t, bool(*core.EvalPure(core.PApp(core.Equal, vs[0], vs[1])).(*core.BooleanType)))
	}
}

func BenchmarkFilter(b *testing.B) {
	go systemt.RunDaemons()

	f := builtinsEnvironment().get("$filter")
	l := randomNumberList(10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		core.EvalPure(core.PApp(core.ToString, core.PApp(f, core.IsOrdered, l)))
	}
}

func BenchmarkFilterBare(b *testing.B) {
	go systemt.RunDaemons()

	f := core.PApp(builtins.Y, core.NewLazyFunction(
		core.NewSignature([]string{"me", "func", "list"}, "", nil, ""),
		func(vs ...core.Value) core.Value {
			f := vs[1]
			l := vs[2]

			return core.PApp(core.If,
				core.PApp(core.Equal, l, core.EmptyList),
				core.EmptyList,
				core.PApp(core.Prepend,
					core.PApp(f, core.PApp(core.First, l)),
					core.PApp(vs[0], f, core.PApp(core.Rest, l))))
		}))

	l := randomNumberList(10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		core.EvalPure(core.PApp(core.ToString, core.PApp(f, core.IsOrdered, l)))
	}
}

func TestSort(t *testing.T) {
	go systemt.RunDaemons()

	for _, vs := range [][2]core.Value{
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
		v := core.PApp(builtinsEnvironment().get("sort"), vs[0])
		t.Log(core.EvalPure(core.PApp(core.ToString, v)))
		assert.True(t, bool(*core.EvalPure(core.PApp(core.Equal, v, vs[1])).(*core.BooleanType)))
	}
}

func TestSortError(t *testing.T) {
	go systemt.RunDaemons()

	_, ok := core.EvalPure(core.App(
		builtinsEnvironment().get("$sort"),
		core.NewArguments(
			[]core.PositionalArgument{
				core.NewPositionalArgument(core.NewList(core.NewNumber(42)), false),
			},
			[]core.KeywordArgument{
				core.NewKeywordArgument("less", builtins.LessEq),
			}))).(*core.ErrorType)

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
		core.EvalPure(core.PApp(core.First, core.PApp(f, l)))
	}
}

func BenchmarkMap(b *testing.B) {
	benchmarkMap(b.N, b.ResetTimer, b.Fail)
}

func benchmarkMap(N int, startTimer, fail func()) {
	go systemt.RunDaemons()

	v := core.PApp(builtinsEnvironment().get("$map"), identity, many42())

	startTimer()

	for i := 0; i < N; i++ {
		if core.NumberType(42) != *core.EvalPure(core.PApp(core.First, v)).(*core.NumberType) {
			fail()
		}

		v = core.PApp(core.Rest, v)
	}
}

func BenchmarkGoMap(b *testing.B) {
	v := many42()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if core.NumberType(42) != *core.EvalPure(core.PApp(core.First, v)).(*core.NumberType) {
			b.Fail()
		}

		v = core.PApp(core.Rest, v)
	}
}

var identity = core.NewLazyFunction(
	core.NewSignature([]string{"arg"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		return vs[0]
	})

func many42() core.Value {
	return core.PApp(core.PApp(builtins.Y, core.NewLazyFunction(
		core.NewSignature([]string{"me"}, "", nil, ""),
		func(vs ...core.Value) core.Value {
			return core.PApp(core.Prepend, core.NewNumber(42), core.PApp(vs[0]))
		})))
}

func randomNumberList(n int) core.Value {
	r := rand.New(rand.NewSource(42))
	vs := make([]core.Value, n)

	for i := range vs {
		vs[i] = core.NewNumber(r.Float64())
	}

	l := core.NewList(vs...)
	core.EvalPure(l)
	return l
}
