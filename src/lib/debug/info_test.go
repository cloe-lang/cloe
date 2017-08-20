package debug

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGoInfo(t *testing.T) {
	assert.Equal(t, "info_test.go", filepath.Base(NewGoInfo(0).file))
}

func TestNewGoInfoWithInvalidSkip(t *testing.T) {
	defer func() {
		assert.NotEqual(t, nil, recover())
	}()

	NewGoInfo(10)
}

func TestLines(t *testing.T) {
	t.Log(NewGoInfo(0).Lines())
}

func TestLinesWithLinePosition(t *testing.T) {
	t.Log(NewInfo("<none>", 1, 1, "(write (+ 123 456))").Lines())
}
