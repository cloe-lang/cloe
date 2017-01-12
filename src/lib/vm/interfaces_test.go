package vm

func testEqual(ts ...*Thunk) bool {
	return bool(App(Equal, ts...).Eval().(boolType))
}
