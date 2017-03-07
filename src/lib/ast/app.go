package ast

import "github.com/raviqqe/tisp/src/lib/debug"

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

func (a App) Function() interface{} {
	return a.function
}

func (a App) Arguments() Arguments {
	return a.args
}

func (a App) DebugInfo() *debug.Info {
	return a.info
}
