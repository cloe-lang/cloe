package http

import "github.com/coel-lang/coel/src/lib/core"

// Module is a module in the language.
var Module = map[string]*core.Thunk{
	"get":         get,
	"getRequests": getRequests,
	"post":        post,
}
