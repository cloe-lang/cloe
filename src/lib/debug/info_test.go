package debug

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGoInfo(t *testing.T) {
	assert.Equal(t, "info_test.go", filepath.Base(NewGoInfo(0).file))
}
