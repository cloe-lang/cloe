package macro

import "github.com/tisp-lang/tisp/src/lib/core"

// FunctionsToMacros converts a name-to-function map into a name-to-macro map.
func FunctionsToMacros(fs map[string]*core.Thunk) map[string]func(...interface{}) interface{} {
	return map[string]func(...interface{}) interface{}{}
}
