// +build performance

package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerformanceMap(t *testing.T) {
	r := testing.Benchmark(BenchmarkMap)
	s := testing.Benchmark(BenchmarkGoMap)
	t.Log(r)
	t.Log(s)
	assert.True(t, r.NsPerOp() < 15*s.NsPerOp())
}
