package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainModule(t *testing.T) {
	for _, s := range []string{"", "(let x 42) (def (f x) (+ x 123)) (write 123)"} {
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
	for _, s := range []string{"", "(let x 123) (def (f x) (+ x 123))"} {
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
	assert.Nil(t, err)
}

func checkFailedResult(t *testing.T, x interface{}, err error) {
	t.Log("Result:", x)
	t.Log("Error:", err)
	assert.NotNil(t, err)
}

func TestImportModule(t *testing.T) {
	for _, str := range []string{`(import "foo")`, `(import "foo/bar")`} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.importModule())()
		assert.Nil(t, err)
	}
}

func TestImportModuleFail(t *testing.T) {
	for _, str := range []string{"(import)", "(import foo)", `(import "\a\b\c\d")`} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.importModule())()
		assert.NotNil(t, err)
	}
}

func TestLetVar(t *testing.T) {
	for _, str := range []string{"(let foo 123)", "(let foo (f x y))"} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.letVar())()
		assert.Nil(t, err)
	}
}

func TestLetMatch(t *testing.T) {
	for _, str := range []string{
		`(let [123 x 789] [123 456 789])`,
		`(let {"foo" x ..rest} {"foo" 42 "bar" 2049})`,
	} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.letMatch())()
		assert.Nil(t, err)
	}
}

func TestDefFunction(t *testing.T) {
	for _, str := range []string{
		"(def (foo) 123)",
		"(def (foo x) (f x y))",
		"(def (foo x y (z 123) (v 456) ..args . a b (c 123) (d 456) ..kwargs) 123)",
	} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.letFunction())()
		assert.Nil(t, err)
	}
}

func TestMutuallyRecursiveDefFunctions(t *testing.T) {
	for _, str := range []string{
		`(mr
			(def (even? n) (if (= n 0) true (odd? (- n 1))))
			(def (odd? n) (if (= n 0) false (even? (- n 1)))))`,
	} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.mutuallyRecursiveDefFunctions())()
		assert.Nil(t, err)
	}
}

func TestSignature(t *testing.T) {
	for _, str := range []string{"", "x", "x y", "(x 123)", "..args", ". x", ". (x 123)", ". ..kwargs", "..args . ..kwargs"} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.signature())()
		assert.Nil(t, err)
	}
}

func TestEffect(t *testing.T) {
	for _, str := range []string{"effect", "..effects", "(foo bar)", "..(foo bar)"} {
		s := newStateWithoutFile(str)
		_, err := s.exhaust(s.effect())()
		assert.Nil(t, err)
	}
}

func TestStringLiteral(t *testing.T) {
	for _, str := range []string{
		`""`,
		`"sl"`,
		"\"   string literal  \n \"",
		`"\""`,
		`"\\"`,
		`"\n"`,
		`"\t"`,
	} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.stringLiteral())()

		t.Logf("%#v", result)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestStrip(t *testing.T) {
	s := newStateWithoutFile("ident  \t ")
	result, err := s.exhaust(s.strip(s.identifier()))()

	t.Logf("%#v", result)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestList(t *testing.T) {
	for _, str := range []string{"[]", "[123 456]", "[(f x) 123]"} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotNil(t, result)
		assert.Nil(t, err)
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
		result, err := s.exhaust(s.expression())()

		t.Logf("%#v", result)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestAnonymousFunction(t *testing.T) {
	for _, str := range []string{
		`(\ (x) x)`,
	} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.anonymousFunction())()

		t.Logf("%#v", result)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestMatchExpression(t *testing.T) {
	for _, str := range []string{
		"(match 123 123 true)",
		"(match (foo bar) [123 ..elems] (process elems) xs (write xs))",
		"(match (foo bar) [\"foo\" 123 ..rest] (process rest) xs (write xs))",
	} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.match())()

		t.Logf("%#v", result)

		assert.NotNil(t, result)
		assert.Nil(t, err)
	}
}

func TestApp(t *testing.T) {
	for _, str := range []string{
		"(f)", "(f x)", "(f x y)", "(f ..x)", "(f . x 123)", "(f . x 123 y 456)",
		"(func . ..kwargs)", "(f ..x (func x y) 123 456 ..foo . a 123 b 456 ..c ..(d 123 456 789))"} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.app())()
		t.Logf("%#v", result)
		assert.Nil(t, err)
	}
}

func TestArguments(t *testing.T) {
	for _, str := range []string{"", "x", "x y", "..x", ". x 123", ". x 123 y 456", ". ..kwargs", "..x (func x y) 123 456 ..foo . a 123 b 456 ..c ..(d 123 456 789)"} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.arguments())()
		t.Logf("%#v", result)
		assert.Nil(t, err)
	}
}

func TestIdentifier(t *testing.T) {
	result, err := newStateWithoutFile(";ident").identifier()()

	t.Log(err)

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func TestIdentifierFail(t *testing.T) {
	for _, str := range []string{"", ".", "..", ".foo", "let"} {
		s := newStateWithoutFile(str)
		result, err := s.identifier()()
		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}

func TestBlank(t *testing.T) {
	for _, str := range []string{"", "   ", "\t", "\n\n", " ; laskdjf \n \t "} {
		s := newStateWithoutFile(str)
		result, err := s.exhaust(s.blank())()

		t.Log(result, err)

		assert.Nil(t, result)
		assert.Nil(t, err)
	}
}

func newStateWithoutFile(source string) *state {
	return newState("", source)
}
