package modules

import (
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/modules/http"
	"github.com/coel-lang/coel/src/lib/modules/json"
)

// Modules is a set of built-in modules
var Modules = map[string]map[string]*core.Thunk{
	"http": http.Module,
	"json": json.Module,
}
