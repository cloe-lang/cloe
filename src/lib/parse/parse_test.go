package parse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModule(t *testing.T) {
	for _, str := range []string{"", "()", "(foo bar)"} {
		result, err := newState(str).module()()

		t.Log(result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestXFailModule(t *testing.T) {
	for _, str := range []string{"(", "(()"} {
		result, err := newState(str).module()()

		t.Log(err.Error())

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range []string{`""`, `"sl"`, "\"   string literal  \n \"", `"\""`, `"\\"`} {
		s := newState(str)
		result, err := s.Exhaust(s.stringLiteral())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestStrip(t *testing.T) {
	s := newState("ident  \t ")
	result, err := s.Exhaust(s.strip(s.identifier()))()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestList(t *testing.T) {
	s := newState("()")
	result, err := s.Exhaust(s.expression())()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestExpression(t *testing.T) {
	strs := []string{
		"ident",
		"ident  ",
		"(foo ; (this is) comment \n bar)  \t ; lsdfj\n ",
	}

	for _, str := range strs {
		t.Logf("source: %#v", str)

		s := newState(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestQuotedExpression(t *testing.T) {
	for _, str := range []string{"`ident", "``ident", "```ident"} {
		s := newState(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestSetLiteral(t *testing.T) {
	s := newState("'{1 2 3}")
	result, err := s.Exhaust(s.expression())()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestClosureLiteral(t *testing.T) {
	s := newState("'(+ #1 #2 3)")
	result, err := s.Exhaust(s.expression())()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestIdentifier(t *testing.T) {
	result, err := newState(";ident").identifier()()

	t.Log(err)

	assert.Equal(t, result, nil)
	assert.NotEqual(t, err, nil)
}

func TestBlank(t *testing.T) {
	for _, str := range []string{"", "   ", "\t", "\n\n", " ; laskdjf \n \t "} {
		s := newState(str)
		result, err := s.Exhaust(s.blank())()

		t.Log(result, err)

		assert.Equal(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestQuote(t *testing.T) {
	for _, str := range []string{"`foo", "`( foo ; lajdfs\n   bar )"} {
		s := newState(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}
