package ast

type App struct {
	function interface{}
	args     Arguments
}

func NewApp(f interface{}, args Arguments) App {
	return App{f, args}
}
