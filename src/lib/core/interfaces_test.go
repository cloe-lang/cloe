package core

import "testing"

func testEqual(t1, t2 *Thunk) bool {
	return compare(t1.Eval(), t2.Eval()) == 0
}

func testLess(t1, t2 *Thunk) bool {
	return t1.Eval().(ordered).compare(t2.Eval().(ordered)) < 0
}

func TestXFailLess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	compare(NewNumber(42).Eval(), NewError("you", "failed.").Eval())
}

func testCompare(t1, t2 *Thunk) NumberType {
	return PApp(Compare, t1, t2).Eval().(NumberType)
}
