package parse

import "testing"

func TestModule1(t *testing.T) {
	for _, str := range []string{"", "()", "(foo bar)"} {
		result, err := newState(str).module()()

		if err == nil {
			t.Log(result)
		} else {
			t.Error(err.Error())
		}
	}
}

func TestXFailModule(t *testing.T) {
	for _, str := range []string{"(", "(()"} {
		result, err := newState(str).module()()

		if err == nil {
			t.Error(result)
		} else {
			t.Log(err.Error())
		}
	}
}

func TestAtom(t *testing.T) {
	s := newState("ident")
	result, err := s.Exhaust(s.atom())()

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%#v", result)
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range []string{`""`, `"sl"`, "\"   string literal  \n \""} {
		s := newState(str)
		result, err := s.Exhaust(s.stringLiteral())()

		if err != nil {
			t.Error(err.Error())
		} else {
			t.Logf("%#v", result)
		}
	}
}

func TestStrip(t *testing.T) {
	s := newState("  ident  ")
	result, err := s.Exhaust(s.strip(s.atom()))()

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%#v", result)
	}
}

func TestWrapChars(t *testing.T) {
	s := newState("( \tident \n)  ; laskdfjsl \t  dkjf\n ")
	result, err := s.Exhaust(s.wrapChars('(', s.atom(), ')'))()

	if err != nil {
		t.Error(err.Error())
	} else {
		t.Logf("%#v", result)
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

func TestElem(t *testing.T) {
	strs := []string{
		"ident",
		"  ident  ",
		" (foo ; (this is) comment \n bar)  \t ; lsdfj\n ",
	}

	for _, str := range strs {
		t.Logf("source: %#v", str)

		s := newState(str)
		result, err := s.Exhaust(s.elem())()

		if err == nil {
			t.Logf("%#v", result)
		} else {
			t.Error(err.Error())
		}
	}
}

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
