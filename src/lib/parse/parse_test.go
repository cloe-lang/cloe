package parse

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainModule(t *testing.T) {
	for _, str := range []string{"", "(let x 42) (let (f x) (+ x 123)) (write 123)"} {
		result, err := newStateWithoutFile(str).mainModule()()

		t.Log(result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestMainModuleWithFiles(t *testing.T) {
	fs, err := filepath.Glob("../../../test/*.tisp")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(fs))

	for _, f := range fs {
		bs, err := ioutil.ReadFile(f)
		assert.Equal(t, nil, err)

		result, err := MainModule(f, string(bs))

		t.Log(result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestXFailMainModule(t *testing.T) {
	for _, str := range []string{"(", "(()"} {
		result, err := newStateWithoutFile(str).mainModule()()

		t.Log(err.Error())

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestSubModule(t *testing.T) {
	for _, str := range []string{"", "(let x 123) (let (f x) (+ x 123))"} {
		result, err := newStateWithoutFile(str).subModule()()

		t.Log(result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestXFailSubModule(t *testing.T) {
	for _, str := range []string{"(", "(()", "(write 123)"} {
		result, err := newStateWithoutFile(str).subModule()()

		t.Log(err.Error())

		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestImportModule(t *testing.T) {
	for _, str := range []string{`(import "foo")`, `(import "foo/bar")`} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.importModule())()
		assert.Equal(t, nil, err)
	}
}

func TestLetVar(t *testing.T) {
	for _, str := range []string{"(let foo 123)", "(let foo (f x y))"} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.letVar())()
		assert.Equal(t, nil, err)
	}
}

func TestLetFunction(t *testing.T) {
	for _, str := range []string{
		"(let (foo) 123)",
		"(let (foo x) (f x y))",
		"(let (foo x y (z 123) (v 456) ..args . a b (c 123) (d 456) ..kwargs) 123)",
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.letFunction())()
		assert.Equal(t, nil, err)
	}
}

func TestSignature(t *testing.T) {
	for _, str := range []string{"", "x", "x y", "(x 123)", "..args", ". x", ". (x 123)", ". ..kwargs", "..args . ..kwargs"} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.signature())()
		assert.Equal(t, nil, err)
	}
}

func TestOutput(t *testing.T) {
	for _, str := range []string{"output", "..outputs", "(foo bar)", "..(foo bar)"} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.output())()
		assert.Equal(t, nil, err)
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range []string{`""`, `"sl"`, "\"   string literal  \n \"", `"\""`, `"\\"`} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.stringLiteral())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestStrip(t *testing.T) {
	s := newStateWithoutFile("ident  \t ")
	result, err := s.Exhaust(s.strip(s.identifier()))()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

func TestList(t *testing.T) {
	for _, str := range []string{"[]", "[123 456]", "[(f x) 123]"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestExpression(t *testing.T) {
	strs := []string{
		"ident",
		"ident  ",
		"(foo ; (this is) comment \n bar)  \t ; lsdfj\n ",
	}

	for _, str := range strs {
		t.Logf("source: %#v", str)

		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestQuotedExpression(t *testing.T) {
	for _, str := range []string{"`ident", "``ident", "```ident"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestSetLiteral(t *testing.T) {
	s := newStateWithoutFile("'{1 2 3}")
	result, err := s.Exhaust(s.expression())()

	t.Logf("%#v", result)

	assert.NotEqual(t, result, nil)
	assert.Equal(t, err, nil)
}

// func TestClosureLiteral(t *testing.T) {
//	s := newStateWithoutFile("'(+ #1 #2 3)")
//	result, err := s.Exhaust(s.expression())()

//	t.Logf("%#v", result)

//	assert.NotEqual(t, result, nil)
//	assert.Equal(t, err, nil)
// }

func TestApp(t *testing.T) {
	for _, str := range []string{
		"(f)", "(f x)", "(f x y)", "(f ..x)", "(f . x 123)", "(f . x 123 y 456)",
		"(func . ..kwargs)", "(f ..x (func x y) 123 456 ..foo . a 123 b 456 ..c ..(d 123 456 789))"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.app())()
		t.Logf("%#v", result)
		assert.Equal(t, err, nil)
	}
}

func TestArguments(t *testing.T) {
	for _, str := range []string{"", "x", "x y", "..x", ". x 123", ". x 123 y 456", ". ..kwargs", "..x (func x y) 123 456 ..foo . a 123 b 456 ..c ..(d 123 456 789)"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.arguments())()
		t.Logf("%#v", result)
		assert.Equal(t, err, nil)
	}
}

func TestIdentifier(t *testing.T) {
	result, err := newStateWithoutFile(";ident").identifier()()

	t.Log(err)

	assert.Equal(t, result, nil)
	assert.NotEqual(t, err, nil)
}

func TestXFailIdentifier(t *testing.T) {
	for _, str := range []string{"", ".", "..", ".foo"} {
		s := newStateWithoutFile(str)
		result, err := s.identifier()()
		assert.Equal(t, result, nil)
		assert.NotEqual(t, err, nil)
	}
}

func TestBlank(t *testing.T) {
	for _, str := range []string{"", "   ", "\t", "\n\n", " ; laskdjf \n \t "} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.blank())()

		t.Log(result, err)

		assert.Equal(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func TestQuote(t *testing.T) {
	for _, str := range []string{"`foo", "`( foo ; lajdfs\n   bar )"} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

func newStateWithoutFile(source string) *state {
	return newState("", source)
}
