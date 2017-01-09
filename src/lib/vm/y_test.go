package vm

import "testing"

func TestY(t *testing.T) {
	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 100} {
		n1 := float64(Y(Normal(NewLazyFunction(lazyFactorial))).(Callable).Call(NumberThunk(n)).(Number))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		if n1 != n2 {
			t.Fail()
		}
	}
}

func strictFactorial(n float64) float64 {
	if n == 0 {
		return 1
	}

	return n * strictFactorial(n-1)
}

func lazyFactorial(ts ...*Thunk) Object {
	return If(
		App(Normal(Equal), ts[1], NumberThunk(0)),
		NumberThunk(1),
		App(Normal(Mult),
			ts[1],
			App(ts[0], App(Normal(Sub), ts[1], NumberThunk(1)))))
}
