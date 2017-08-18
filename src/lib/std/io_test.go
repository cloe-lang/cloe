package std

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tisp-lang/tisp/src/lib/core"
)

func TestReadWithStdin(t *testing.T) {
	assert.Equal(t, core.StringType(""), core.PApp(Read).Eval())
}
