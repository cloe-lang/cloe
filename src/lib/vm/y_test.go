package vm

import "testing"

func TestY(t *testing.T) {
	for _, n := range []float64{0, 1, 2, 3, 4, 5, 6, 100} {
		n1 := lazyFactorial(NumberThunk(n))
		n2 := strictFactorial(n)

		t.Logf("%d: %f == %f?\n", int(n), n1, n2)

		if n1 != n2 {
			t.Fail()
		}
	}

	for _, ts := range [][]*Thunk{
		{NumberThunk(7)},
		{NumberThunk(13), StringThunk("foobarbaz")},
		{NumberThunk(42), Nil, Nil},
	} {
		t.Log(lazyFactorial(ts...))
	}
}

func strictFactorial(n float64) float64 {
	if n == 0 {
		return 1
	}

	return n * strictFactorial(n-1)
}

func lazyFactorial(ts ...*Thunk) float64 {
	return float64(App(App(Y, NewLazyFunction(lazyFactorialImpl)), ts...).Eval().(Number))
}

func lazyFactorialImpl(ts ...*Thunk) Object {
	// fmt.Println(len(ts))

	return App(If,
		App(Equal, ts[1], NumberThunk(0)),
		NumberThunk(1),
		App(Mult,
			ts[1],
			App(ts[0], append([]*Thunk{App(Sub, ts[1], NumberThunk(1))}, ts[2:]...)...)))
}
