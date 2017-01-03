package comb

import "testing"

func TestMany(t *testing.T) {
	for _, str := range []string{"", "  "} {
		s := NewState(str)
		result, err := s.Many(s.Char(' '))()

		if err != nil {
			t.Error(err.Error())
		}

		t.Logf("%#v", result)
	}
}

func testMany1Space(str string) (interface{}, error) {
	s := NewState(str)
	return s.Many1(s.Char(' '))()
}

func TestMany1(t *testing.T) {
	result, err := testMany1Space(" ")

	if err != nil {
		t.Error(err.Error())
	}

	t.Logf("%#v", result)
}

func TestMany1Fail(t *testing.T) {
	result, err := testMany1Space("")

	if result != nil {
		t.Errorf("`result` should be nil but %#v.", result)
	}

	t.Log(err.Error())
}
