package ast

type App struct {
	function interface{}
	args     Arguments
}

func NewApp(f interface{}, args Arguments) App {
	return App{f, args}
}

func (a App) Function() interface{} {
	return a.function
}

func (a App) Arguments() Arguments {
	return a.args
}
