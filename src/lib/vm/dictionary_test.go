package vm

import "testing"

func TestDictionarySet(t *testing.T) {
	for _, th := range []*Thunk{
		True, False, Nil, NewNumber(42), NewString("risp"),
	} {
		_, ok := EmptyDictionary.Eval().(dictionaryType).Set(th.Eval(), Nil).(dictionaryType)

		if !ok {
			t.Fail()
		}
	}
}
