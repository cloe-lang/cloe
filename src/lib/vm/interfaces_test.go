package vm

func testEqual(ts ...*Thunk) bool {
	return bool(App(Normal(Equal), ts...).EvalStrictly().(Bool))
}
