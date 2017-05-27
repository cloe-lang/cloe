package ast

import (
	"github.com/tisp-lang/tisp/src/lib/debug"
)

// App represents an application of a function to arguments.
type App struct {
	function interface{}
	args     Arguments
	info     debug.Info
}

// NewApp creates an App from a function and arguments with debug information.
func NewApp(f interface{}, args Arguments, info debug.Info) App {
	return App{f, args, info}
}

// NewPApp don't creates PPap but PApp.
func NewPApp(f interface{}, args []interface{}, info debug.Info) App {
	ps := make([]PositionalArgument, 0, len(args))

	for _, arg := range args {
		ps = append(ps, NewPositionalArgument(arg, false))
	}

	return App{f, NewArguments(ps, nil, nil), info}
}

// Function returns a function of an application.
func (a App) Function() interface{} {
	return a.function
}

// Arguments returns arguments of an application.
func (a App) Arguments() Arguments {
	return a.args
}

// DebugInfo returns debug information of an application.
func (a App) DebugInfo() debug.Info {
	return a.info
}
