package vm

func testEqual(ts ...*Thunk) bool {
	return bool(App(Normal(Equal), ts...).Eval().(boolType))
}
