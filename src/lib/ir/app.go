package ir

import (
	"github.com/tisp-lang/tisp/src/lib/core"
	"github.com/tisp-lang/tisp/src/lib/debug"
)

// App represents an application of a function to arguments.
type App struct {
	function interface{}
	args     Arguments
	info     debug.Info
}

// NewApp creates an App from a function and arguments of expressions in IR.
func NewApp(f interface{}, args Arguments, info debug.Info) App {
	return App{f, args, info}
}

func (app App) interpret(args []*core.Thunk) *core.Thunk {
	ps := make([]core.PositionalArgument, len(app.args.positionals))

	for i, p := range app.args.positionals {
		ps[i] = p.interpret(args)
	}

	ks := make([]core.KeywordArgument, len(app.args.keywords))

	for i, k := range app.args.keywords {
		ks[i] = k.interpret(args)
	}

	ds := make([]*core.Thunk, len(app.args.expandedDicts))

	for i, d := range app.args.expandedDicts {
		ds[i] = interpretExpression(args, d)
	}

	return core.AppWithInfo(interpretExpression(args, app.function), core.NewArguments(ps, ks, ds), app.info)
}
