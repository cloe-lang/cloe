package ir

import (
	"github.com/cloe-lang/cloe/src/lib/core"
)

const (
	expandedKeywordArgument = 255
	switchExpression        = 255
)

// Interpreter represents an IR byte code interpreter.
type Interpreter struct {
	code      []byte
	switches  []switchData
	names     []string
	index     int
	variables []core.Value
}

// NewInterpreter creates a new interpreter.
func NewInterpreter(bs []byte, ss []switchData, ns []string, vs []core.Value) Interpreter {
	return Interpreter{bs, ss, ns, 0, vs}
}

// Interpret interprets byte code and returns a result value.
func (j *Interpreter) Interpret() core.Value {
	for j.index < len(j.code) {
		var x core.Value

		switch c := j.readCode(); c {
		case switchExpression:
			v := j.getVariable()
			s := j.switches[j.readCode()]
			dc := j.getVariableWithIndex(s.DefaultCase)

			kvs := make([]core.KeyValue, 0, len(s.Cases))

			for _, c := range s.Cases {
				kvs = append(kvs, core.KeyValue{
					Key:   c.Value,
					Value: j.getVariableWithIndex(c.VariableIndex),
				})
			}

			d := core.NewDictionary(kvs)

			x = core.PApp(core.NewLazyFunction(
				core.NewSignature(nil, "", nil, ""),
				func(...core.Value) core.Value {
					b, err := core.EvalBoolean(core.PApp(core.Include, d, v))

					if err != nil {
						// TODO: Distinguish matching errors and `return err`.
						return dc
					} else if b {
						dc = core.PApp(core.Index, d, v)
					}

					return dc
				}))
		default:
			x = core.App(j.getVariableWithIndex(int(c)), j.getArguments())
		}

		if b := j.readCode(); int(b) < len(j.variables) {
			j.variables[b] = x
		} else {
			j.variables = append(j.variables, x)
		}
	}

	return j.getVariableWithIndex(0)
}

func (j *Interpreter) getArguments() core.Arguments {
	return core.NewArguments(j.getPositionalArguments(), j.getKeywordArguments())
}

func (j *Interpreter) getPositionalArguments() []core.PositionalArgument {
	ps := []core.PositionalArgument{}
	n := j.readCode()

	for i := 0; i < int(n); i++ {
		b := false

		if j.readCode() != 0 {
			b = true
		}

		ps = append(ps, core.NewPositionalArgument(j.getVariable(), b))
	}

	return ps
}

func (j *Interpreter) getKeywordArguments() []core.KeywordArgument {
	ks := []core.KeywordArgument{}
	n := j.readCode()

	for i := 0; i < int(n); i++ {
		if b := j.readCode(); int(b) == expandedKeywordArgument {
			ks = append(ks, core.NewKeywordArgument("", j.getVariable()))
		} else {
			ks = append(ks, core.NewKeywordArgument(j.names[b], j.getVariable()))
		}
	}

	return ks
}

func (j *Interpreter) getVariable() core.Value {
	return j.variables[j.readCode()]
}

func (j *Interpreter) getVariableWithIndex(i int) core.Value {
	return j.variables[i]
}

func (j *Interpreter) readCode() byte {
	b := j.code[j.index]
	j.index++
	return b
}
