package ir

import (
	"github.com/coel-lang/coel/src/lib/core"
	"github.com/coel-lang/coel/src/lib/debug"
)

// App represents an application of a function to arguments.
type App struct {
	function interface{}
	args     Arguments
	info     *debug.Info
}

// NewApp creates an App from a function and arguments of expressions in IR.
func NewApp(f interface{}, args Arguments, info *debug.Info) App {
	return App{f, args, info}
}

func (app App) interpret(args []core.Value) core.Value {
	ps := make([]core.PositionalArgument, 0, len(app.args.positionals))

	for _, p := range app.args.positionals {
		ps = append(ps, p.interpret(args))
	}

	ks := make([]core.KeywordArgument, 0, len(app.args.keywords))

	for _, k := range app.args.keywords {
		ks = append(ks, k.interpret(args))
	}

	return core.AppWithInfo(interpretExpression(args, app.function), core.NewArguments(ps, ks), app.info)
}
