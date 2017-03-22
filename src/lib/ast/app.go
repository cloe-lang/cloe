package ast

import (
	"github.com/raviqqe/tisp/src/lib/debug"
)

type App struct {
	function interface{}
	args     Arguments
	info     debug.Info
}

func NewApp(f interface{}, args Arguments, info debug.Info) App {
	return App{f, args, info}
}

func NewPApp(f interface{}, args []interface{}, info debug.Info) App {
	ps := make([]PositionalArgument, 0, len(args))

	for _, arg := range args {
		ps = append(ps, NewPositionalArgument(arg, false))
	}

	return App{f, NewArguments(ps, nil, []interface{}{}), info}
}

func (a App) Function() interface{} {
	return a.function
}

func (a App) Arguments() Arguments {
	return a.args
}

func (a App) DebugInfo() debug.Info {
	return a.info
}
