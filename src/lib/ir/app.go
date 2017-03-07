package ir

import (
	"github.com/raviqqe/tisp/src/lib/core"
	"github.com/raviqqe/tisp/src/lib/debug"
)

type App struct {
	function interface{}
	args     Arguments
	info     *debug.Info
}

func NewApp(f interface{}, args Arguments) App {
	return App{f, args, nil}
}

func NewAppWithInfo(f interface{}, args Arguments, info *debug.Info) App {
	return App{f, args, info}
}

func (app App) compile(args []*core.Thunk) *core.Thunk {
	ps := make([]core.PositionalArgument, len(app.args.positionals))

	for i, p := range app.args.positionals {
		ps[i] = p.compile(args)
	}

	ks := make([]core.KeywordArgument, len(app.args.keywords))

	for i, k := range app.args.keywords {
		ks[i] = k.compile(args)
	}

	ds := make([]*core.Thunk, len(app.args.expandedDicts))

	for i, d := range app.args.expandedDicts {
		ds[i] = compileExpression(args, d)
	}

	return core.AppWithInfo(compileExpression(args, app.function), core.NewArguments(ps, ks, ds), app.info)
}
