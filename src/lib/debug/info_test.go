package debug

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGoInfo(t *testing.T) {
	Debug = true
	assert.Equal(t, "info_test.go", filepath.Base(NewGoInfo(0).file))
}

func TestNewGoInfoWithInvalidSkip(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	Debug = true
	NewGoInfo(10)
}

func TestLines(t *testing.T) {
	Debug = true
	t.Log(NewGoInfo(0).Lines())
}

func TestLinesWithLinePosition(t *testing.T) {
	Debug = true
	t.Log(NewInfo("<none>", 1, 1, "(print (+ 123 456))").Lines())
}

func TestInfoLinesEmpty(t *testing.T) {
	Debug = false
	assert.Equal(t, "", NewGoInfo(0).Lines())
}
