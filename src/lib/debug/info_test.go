package debug

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoInfo(t *testing.T) {
	assert.Equal(t, "info_test.go", filepath.Base(NewGoInfo(0).file))
}
