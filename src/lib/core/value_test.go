package core

import "testing"

func BenchmarkValueCopyBoolean(b *testing.B) {
	generateCopyBenchmark(True)(b)
}

func BenchmarkValueCopyDictionary(b *testing.B) {
	generateCopyBenchmark(NewDictionary([]KeyValue{{NewString("foo"), NewNumber(42)}}))(b)
}

func BenchmarkValueCopyFunction(b *testing.B) {
	generateCopyBenchmark(If)(b)
}

func BenchmarkValueCopyList(b *testing.B) {
	generateCopyBenchmark(NewList(NewNumber(42)))(b)
}

func BenchmarkValueCopyNil(b *testing.B) {
	generateCopyBenchmark(Nil)(b)
}

func BenchmarkValueCopyNumber(b *testing.B) {
	generateCopyBenchmark(NewNumber(42))(b)
}

func BenchmarkValueCopyString(b *testing.B) {
	generateCopyBenchmark(NewString("foo"))(b)
}

func BenchmarkValueCopyError(b *testing.B) {
	generateCopyBenchmark(NewError("MyError", "No way!"))(b)
}

type UnboxedNumberType float64

func (b UnboxedNumberType) eval() Value {
	return b
}

func BenchmarkValueCopyInefficientUnboxedType(b *testing.B) {
	// The result must be 1 allocs/op.
	generateCopyBenchmark(UnboxedNumberType(42))(b)
}

func generateCopyBenchmark(v Value) func(b *testing.B) {
	return func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			v = v.eval()
		}
	}
}
