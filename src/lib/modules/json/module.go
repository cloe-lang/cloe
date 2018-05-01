package json

import "github.com/cloe-lang/cloe/src/lib/core"

// Module is a module in the language.
var Module = map[string]core.Value{
	"decode": decode,
	"encode": encode,
}
