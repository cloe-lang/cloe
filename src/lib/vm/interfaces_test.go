package vm

func testEqual(ts ...*Thunk) bool {
	return bool(Equal(ts...).Eval().(Bool))
}
