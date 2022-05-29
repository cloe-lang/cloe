//go:build performance

package compile

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerformanceMap(t *testing.T) {
	r := testing.Benchmark(BenchmarkMap)
	s := testing.Benchmark(BenchmarkGoMap)
	t.Log(r)
	t.Log(s)
	assert.True(t, r.NsPerOp() < 15*s.NsPerOp())
}

func TestPerformanceMapOrder(t *testing.T) {
	b := func(N int) float64 {
		var start time.Time
		benchmarkMap(N, func() { start = time.Now() }, t.Fail)
		return time.Since(start).Seconds()
	}

	r := b(10000) / b(2000)
	t.Log(r)
	assert.True(t, 4 < r && r < 6)
}
