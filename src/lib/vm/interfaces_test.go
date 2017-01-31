package vm

func testEqual(ts ...*Thunk) bool {
	return bool(App(Equal, ts...).Eval().(boolType))
}

func testLess(t1, t2 *Thunk) bool {
	return t1.Eval().(ordered).less(t2.Eval().(ordered))
}
