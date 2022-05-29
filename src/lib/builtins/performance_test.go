//go:build performance

package builtins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerformanceY(t *testing.T) {
	r := testing.Benchmark(BenchmarkY)
	s := testing.Benchmark(BenchmarkGoY)
	t.Log(r)
	t.Log(s)
	assert.True(t, r.NsPerOp() < 10*s.NsPerOp())
}
