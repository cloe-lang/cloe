package ir

import (
	"testing"

	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestNewInterpreter(t *testing.T) {
	NewInterpreter(nil, nil, nil, nil)
}

func TestInterpreterInterpret(t *testing.T) {
	for _, c := range []struct {
		interpreter Interpreter
		value       core.Value
	}{
		{
			NewInterpreter(nil, nil, nil, []core.Value{core.Nil}),
			core.Nil,
		},
		{
			NewInterpreter([]byte{0, 0, 0, 0}, nil, nil, []core.Value{core.Add}),
			core.NewNumber(0),
		},
		{
			NewInterpreter([]byte{0, 0, 0, 1, 0, 1, 0, 1, 0, 0}, nil, nil, []core.Value{core.Add}),
			core.NewNumber(0),
		},
		{
			NewInterpreter(
				[]byte{0, 1, 0, 1, 0, 0},
				nil,
				nil,
				[]core.Value{core.Add, core.NewNumber(1)}),
			core.NewNumber(1),
		},
		{
			NewInterpreter(
				[]byte{0, 2, 0, 1, 0, 2, 0, 0},
				nil,
				nil,
				[]core.Value{core.Add, core.NewNumber(1), core.NewNumber(2)}),
			core.NewNumber(3),
		},
		{
			NewInterpreter(
				[]byte{0, 2, 0, 1, 0, 2, 0, 0},
				nil,
				nil,
				[]core.Value{core.Add, core.NewNumber(1), core.NewNumber(2)}),
			core.NewNumber(3),
		},
		{
			NewInterpreter(
				[]byte{0, 1, 1, 1, 0, 0},
				nil,
				nil,
				[]core.Value{core.Add, core.NewList(core.NewNumber(1), core.NewNumber(2))}),
			core.NewNumber(3),
		},
		{
			NewInterpreter(
				[]byte{0, 0, 1, 0, 1, 0},
				nil,
				[]string{"foo"},
				[]core.Value{testFunction, core.True}),
			core.True,
		},
		{
			NewInterpreter(
				[]byte{0, 0, 1, expandedKeywordArgument, 1, 0},
				nil,
				nil,
				[]core.Value{
					testFunction,
					core.NewDictionary([]core.KeyValue{{Key: core.NewString("foo"), Value: core.True}}),
				}),
			core.True,
		},
		{
			NewInterpreter(
				[]byte{switchExpression, 1, 0, 0},
				[]switchData{newSwitchData([]caseData{newCaseData(core.NewString("foo"), 2)}, 0)},
				nil,
				[]core.Value{core.DummyError, core.NewString("foo"), core.True}),
			core.True,
		},
	} {
		b, err := core.EvalBoolean(core.PApp(core.Equal, c.value, c.interpreter.Interpret()))

		assert.Nil(t, err)
		assert.True(t, bool(b))
	}
}

func TestInterpreterInterpretError(t *testing.T) {
	for _, c := range []Interpreter{
		NewInterpreter(
			[]byte{switchExpression, 1, 0, 0},
			[]switchData{newSwitchData([]caseData{newCaseData(core.NewString("foo"), 0)}, 0)},
			nil,
			[]core.Value{core.DummyError, core.NewString("foo")}),
		NewInterpreter(
			[]byte{switchExpression, 1, 0, 0},
			[]switchData{newSwitchData([]caseData{newCaseData(core.DummyError, 1)}, 1)},
			nil,
			[]core.Value{core.DummyError, core.NewString("foo")}),
	} {
		_, err := core.EvalBoolean(c.Interpret())

		assert.NotNil(t, err)
	}
}

var testFunction = core.NewLazyFunction(
	core.NewSignature(
		nil, "",
		[]core.OptionalParameter{core.NewOptionalParameter("foo", core.Nil)}, "",
	),
	func(vs ...core.Value) core.Value {
		return vs[0]
	})
