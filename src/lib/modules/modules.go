package modules

import (
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/modules/fs"
	"github.com/coel-lang/coel/src/lib/modules/http"
	"github.com/coel-lang/coel/src/lib/modules/json"
	"github.com/coel-lang/coel/src/lib/modules/re"
)

// Modules is a set of built-in modules
var Modules = map[string]map[string]*core.Thunk{
	"http": http.Module,
	"json": json.Module,
	"fs":   fs.Module,
	"re":   re.Module,
}
