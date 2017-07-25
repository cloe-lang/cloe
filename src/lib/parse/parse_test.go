package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainModule(t *testing.T) {
	for _, s := range []string{"", "(let x 42) (let (f x) (+ x 123)) (write 123)"} {
		ss, err := MainModule("<test>", s)
		checkSuccessfulResult(t, ss, err)
		x, err := newStateWithoutFile(s).mainModule()()
		checkSuccessfulResult(t, x, err)
	}
}

func TestMainModuleFail(t *testing.T) {
	for _, s := range []string{"(", "(()"} {
		ss, err := MainModule("<test>", s)
		checkFailedResult(t, ss, err)
		x, err := newStateWithoutFile(s).mainModule()()
		checkFailedResult(t, x, err)
	}
}

func TestSubModule(t *testing.T) {
	for _, s := range []string{"", "(let x 123) (let (f x) (+ x 123))"} {
		ss, err := SubModule("<test>", s)
		checkSuccessfulResult(t, ss, err)
		x, err := newStateWithoutFile(s).subModule()()
		checkSuccessfulResult(t, x, err)
	}
}

func TestSubModuleFail(t *testing.T) {
	for _, s := range []string{"(", "(()", "(write 123)"} {
		ss, err := SubModule("<test>", s)
		checkFailedResult(t, ss, err)
		x, err := newStateWithoutFile(s).subModule()()
		checkFailedResult(t, x, err)
	}
}

func checkSuccessfulResult(t *testing.T, x interface{}, err error) {
	t.Log("Result:", x)
	t.Log("Error:", err)
	assert.Equal(t, nil, err)
}

func checkFailedResult(t *testing.T, x interface{}, err error) {
	t.Log("Result:", x)
	t.Log("Error:", err)
	assert.NotEqual(t, nil, err)
}

func TestImportModule(t *testing.T) {
	for _, str := range []string{`(import "foo")`, `(import "foo/bar")`} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.importModule())()
		assert.Equal(t, nil, err)
	}
}

func TestImportModuleFail(t *testing.T) {
	for _, str := range []string{"(import)", "(import foo)", `(import "\a\b\c\d")`} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.importModule())()
		assert.NotEqual(t, nil, err)
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

func TestMutuallyRecursiveLetFunctions(t *testing.T) {
	for _, str := range []string{
		`(mr
			(let (even? n) (if (= n 0) true (odd? (- n 1))))
			(let (odd? n) (if (= n 0) false (even? (- n 1)))))`,
	} {
		s := newStateWithoutFile(str)
		_, err := s.Exhaust(s.mutuallyRecursiveLetFunctions())()
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

func TestMatchExpression(t *testing.T) {
	for _, str := range []string{
		"(match 123 123 true)",
		"(match (foo bar) [123 ..elems] (process elems) xs (write xs))",
		"(match (foo bar) [\"foo\" 123 ..rest] (process rest) xs (write xs))",
	} {
		s := newStateWithoutFile(str)
		result, err := s.Exhaust(s.match())()

		t.Logf("%#v", result)

		assert.NotEqual(t, result, nil)
		assert.Equal(t, err, nil)
	}
}

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

func TestIdentifierFail(t *testing.T) {
	for _, str := range []string{"", ".", "..", ".foo", "let"} {
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

func newStateWithoutFile(source string) *state {
	return newState("", source)
}

func TestUnquoteStringFail(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	unquoteString("foo")
}

// func TestClosureLiteral(t *testing.T) {
//	s := newStateWithoutFile("'(+ #1 #2 3)")
//	result, err := s.Exhaust(s.expression())()

//	t.Logf("%#v", result)

//	assert.NotEqual(t, result, nil)
//	assert.Equal(t, err, nil)
// }
