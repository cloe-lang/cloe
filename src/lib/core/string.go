package core

import (
	"fmt"
	"strings"
)

// StringType represents a string in the language.
type StringType string

// Eval evaluates a value into a WHNF.
func (s *StringType) eval() Value {
	return s
}

// NewString creates a string in the language from one in Go.
func NewString(s string) *StringType {
	ss := StringType(s)
	return &ss
}

func (s *StringType) call(args Arguments) Value {
	return Index.call(NewPositionalArguments(s).Merge(args))
}

func (s *StringType) index(v Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	i := int(n)
	rs := []rune(string(*s))

	if i > len(rs) {
		return OutOfRangeError()
	}

	return NewString(string(rs[i-1 : i]))
}

func (s *StringType) insert(v Value, t Value) Value {
	n, err := checkIndex(v)

	if err != nil {
		return err
	}

	i := int(n) - 1

	if i > len(*s) {
		return OutOfRangeError()
	}

	ss, err := EvalString(t)

	if err != nil {
		return err
	}

	*ss = (*s)[:i] + *ss + (*s)[i:]
	return ss
}

func (s *StringType) merge(vs ...Value) Value {
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

		ss = append(ss, string(*s))
	}

	return NewString(strings.Join(ss, ""))
}

func (s *StringType) delete(v Value) Value {
	n, err := checkIndex(v)
	i := int(n) - 1

	if err != nil {
		return err
	} else if i >= len(*s) {
		return OutOfRangeError()
	}

	ss := (*s)[:i] + (*s)[i+1:]
	return &ss
}

func (s *StringType) toList() Value {
	if *s == "" {
		return EmptyList
	}

	rs := []rune(string(*s))

	return cons(
		NewString(string(rs[0])),
		PApp(ToList, NewString(string(rs[1:]))))
}

func (s *StringType) compare(c comparable) int {
	return strings.Compare(string(*s), string(*c.(*StringType)))
}

func (*StringType) ordered() {}

func (s *StringType) string() Value {
	return s
}

func (s *StringType) quoted() Value {
	return NewString(fmt.Sprintf("%#v", string(*s)))
}

func (s *StringType) size() Value {
	return NewNumber(float64(len([]rune(string(*s)))))
}

func (s *StringType) include(v Value) Value {
	ss, ok := v.(*StringType)

	if !ok {
		return NotStringError(v)
	}

	return NewBool(strings.Contains(string(*s), string(*ss)))
}

func (s *StringType) String() string {
	return string(*s)
}
