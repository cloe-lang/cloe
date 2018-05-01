package modules

import (
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/cloe-lang/cloe/src/lib/modules/fs"
	"github.com/cloe-lang/cloe/src/lib/modules/http"
	"github.com/cloe-lang/cloe/src/lib/modules/json"
	"github.com/cloe-lang/cloe/src/lib/modules/os"
	"github.com/cloe-lang/cloe/src/lib/modules/re"
)

// Modules is a set of built-in modules
var Modules = map[string]map[string]core.Value{
	"fs":   fs.Module,
	"http": http.Module,
	"json": json.Module,
	"os":   os.Module,
	"re":   re.Module,
}
