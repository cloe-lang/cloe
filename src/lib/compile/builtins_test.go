package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuiltinsEnvironment(t *testing.T) {
	builtinsEnvironment()
}

func TestCompileBuiltinModule(t *testing.T) {
	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) x)`)
}

func TestCompileBuiltinModuleWithInvalidSyntax(t *testing.T) {
	defer func() {
		_, ok := recover().(error)
		assert.True(t, ok)
	}()

	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) x`)
}

func TestCompileBuiltinModuleWithInvalidSource(t *testing.T) {
	defer func() {
		_, ok := recover().(error)
		assert.True(t, ok)
	}()

	compileBuiltinModule(newEnvironment(testFallback), "", `(def (foo x) y)`)
}
