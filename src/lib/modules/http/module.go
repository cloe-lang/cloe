package http

import "github.com/cloe-lang/cloe/src/lib/core"

// Module is a module in the language.
var Module = map[string]core.Value{
	"get":         get,
	"getRequests": getRequests,
	"post":        post,
}
