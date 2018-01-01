package modules

import (
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/modules/fs"
	"github.com/coel-lang/coel/src/lib/modules/http"
	"github.com/coel-lang/coel/src/lib/modules/json"
	"github.com/coel-lang/coel/src/lib/modules/os"
	"github.com/coel-lang/coel/src/lib/modules/re"
)

// Modules is a set of built-in modules
var Modules = map[string]map[string]*core.Thunk{
	"fs":   fs.Module,
	"http": http.Module,
	"json": json.Module,
	"os":   os.Module,
	"re":   re.Module,
}
