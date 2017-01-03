package parse

import "testing"

// func TestModule(t *testing.T) {
// 	newState("").module()()
// }

func TestAtom(t *testing.T) {
	s := newState("ident")
	result, err := s.atom()()

	if !s.Exhausted() {
		t.Error("Source is not exhausted.")
	}

	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("%#v", toString(result))
}

// func TestElem(t *testing.T) {
// 	result, err := newState("   ident  ").atom()()

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	t.Logf(`"%s"`, toString(result))
// }

func TestBlank(t *testing.T) {
	for _, str := range []string{"", "   ", "\t", "\n\n", " ; laskdjf \n \t "} {
		s := newState(str)
		result, err := s.blank()()

		if !s.Exhausted() {
			t.Error("Source is not exhausted.")
		}

		if result != nil {
			t.Errorf("`result` should be null. (%#v)", result)
		}

		if err != nil {
			t.Errorf("`err` should be null. (%#v)", result)
		}
	}
}

func toString(any interface{}) string {
	xs := any.([]interface{})
	rs := make([]rune, len(xs))

	for i, x := range xs {
		rs[i] = x.(rune)
	}

	return string(rs)
}
