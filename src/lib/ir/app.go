package ir

import "../vm"

type App struct {
	function interface{}
	args     Arguments
}

func NewApp(f interface{}, args Arguments) App {
	return App{f, args}
}

func (app App) compile(args []*vm.Thunk) *vm.Thunk {
	ps := make([]vm.PositionalArgument, len(app.args.positionals))

	for i, p := range app.args.positionals {
		ps[i] = p.compile(args)
	}

	ks := make([]vm.KeywordArgument, len(app.args.keywords))

	for i, k := range app.args.keywords {
		ks[i] = k.compile(args)
	}

	ds := make([]*vm.Thunk, len(app.args.expandedDicts))

	for i, d := range app.args.expandedDicts {
		ds[i] = compileExpression(args, d)
	}

	return vm.App(compileExpression(args, app.function), vm.NewArguments(ps, ks, ds))
}
