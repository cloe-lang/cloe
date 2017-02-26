package core

func testEqual(ts ...*Thunk) bool {
	return bool(PApp(Equal, ts...).Eval().(BoolType))
}

func testLess(t1, t2 *Thunk) bool {
	return t1.Eval().(ordered).less(t2.Eval().(ordered))
}
