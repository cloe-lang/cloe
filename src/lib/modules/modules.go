package modules

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/modules/http"
)

// Modules is a set of built-in modules
var Modules = map[string]map[string]*core.Thunk{
	"http": http.Module,
}
