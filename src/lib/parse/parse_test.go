package parse

import "testing"

// func TestModule(t *testing.T) {
// 	newState("").module()()
// }

func TestAtom(t *testing.T) {
	s := newState("ident")
	result, err := s.Exhaust(s.atom())()

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%#v", toString(result))
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range []string{`""`, `"sl"`, "\"   string literal  \n \""} {
		s := newState(str)
		result, err := s.Exhaust(s.stringLiteral())()

		if err != nil {
			t.Error(err.Error())
		} else {
			t.Logf("%#v", toString(result))
		}
	}
}

func TestStrip(t *testing.T) {
	s := newState("  ident  ")
	result, err := s.Exhaust(s.strip(s.atom()))()

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%#v", toString(result))
	}
}

func TestWrapChars(t *testing.T) {
	s := newState(" ; laskdfjsl \t  dkjf\n ( \tident \n)  ")
	result, err := s.Exhaust(s.wrapChars('(', s.atom(), ')'))()

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%#v", toString(result))
	}
}

func TestList(t *testing.T) {
	s := newState("()")
	result, err := s.Exhaust(s.list())()

	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("%#v", result)
}

// func TestElem(t *testing.T) {
// 	for _, str := range []string{"ident", "  ident  "} {
// 		t.Logf("source: %#v", str)

// 		s := newState(str)
// 		result, err := s.elem()()

// 		if err != nil {
// 			t.Error(err.Error())
// 		} else {
// 			t.Logf("%#v", toString(result))
// 		}
// 	}
// }

func TestBlank(t *testing.T) {
	for _, str := range []string{"", "   ", "\t", "\n\n", " ; laskdjf \n \t "} {
		s := newState(str)
		result, err := s.Exhaust(s.blank())()

		t.Log(result, err)

		if result != nil {
			t.Errorf("`result` should be nil. (%#v)", result)
		}

		if err != nil {
			t.Errorf("`err` should be nil. (%#v)", result)
		}
	}
}
