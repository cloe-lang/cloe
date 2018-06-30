package core

import (
	"fmt"
	"strings"
)

// StringType represents a string in the language.
type StringType string

// Eval evaluates a value into a WHNF.
func (s StringType) eval() Value {
	return s
}

// NewString creates a string in the language from one in Go.
func NewString(s string) StringType {
	return StringType(s)
}

func (s StringType) assign(i, v Value) Value {
	n, err := checkIndex(i)

	if err != nil {
		return err
	}

	rs := []rune(string(s))

	if int(n) > len(rs) {
		return OutOfRangeError()
	}

	s, err = EvalString(v)

	if err != nil {
		return err
	}

	cs := []rune(string(s))

	if len(cs) != 1 {
		return ValueError("cannot assign non-character to string")
	}

	rs[int(n)-1] = cs[0]

	return StringType(rs)
}

func (s StringType) index(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	i := int(n)
	rs := []rune(string(s))

	if i > len(rs) {
		return OutOfRangeError()
	}

	return NewString(string(rs[i-1 : i]))
}

func (s StringType) insert(v Value, t Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	i := int(n) - 1

	if i > len(s) {
		return OutOfRangeError()
	}

	ss, err := EvalString(t)

	if err != nil {
		return err
	}

	return s[:i] + ss + s[i:]
}

func (s StringType) merge(vs ...Value) Value {
	vs = append([]Value{s}, vs...)

	for _, t := range vs {
		go EvalPure(t)
	}

	ss := make([]string, 0, len(vs))

	for _, t := range vs {
		s, err := EvalString(t)

		if err != nil {
			return err
		}

		ss = append(ss, string(s))
	}

	return NewString(strings.Join(ss, ""))
}

func (s StringType) delete(v Value) Value {
	n, err := checkIndex(v)
	i := int(n) - 1

	if err != nil {
		return err
	} else if i >= len(s) {
		return OutOfRangeError()
	}

	return s[:i] + s[i+1:]
}

func (s StringType) toList() Value {
	if s == "" {
		return EmptyList
	}

	rs := []rune(string(s))

	return cons(
		NewString(string(rs[0])),
		PApp(ToList, NewString(string(rs[1:]))))
}

func (s StringType) compare(c comparable) int {
	return strings.Compare(string(s), string(c.(StringType)))
}

func (StringType) ordered() {}

func (s StringType) string() Value {
	return s
}

func (s StringType) quoted() Value {
	return NewString(fmt.Sprintf("%#v", string(s)))
}

func (s StringType) size() Value {
	return NewNumber(float64(len([]rune(string(s)))))
}

func (s StringType) include(v Value) Value {
	ss, ok := v.(StringType)

	if !ok {
		return NotStringError(v)
	}

	return NewBoolean(strings.Contains(string(s), string(ss)))
}
