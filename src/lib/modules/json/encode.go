package json

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coel-lang/coel/src/lib/core"
)

var jsonEncodeError = jsonError(errors.New("Cannot be encoded as JSON"))

var encode = core.NewLazyFunction(
	core.NewSignature([]string{"decoded"}, nil, "", nil, nil, ""),
	func(ts ...*core.Thunk) core.Value {
		s, err := encodeValue(ts[0].Eval())

		if err != nil {
			return err
		}

		if s == "" {
			return jsonEncodeError
		}

		return core.NewString(s)
	})

func encodeValue(v core.Value) (result string, err *core.Thunk) {
	switch v := v.(type) {
	case core.NilType:
		return "null", nil
	case core.NumberType:
		return fmt.Sprintf("%v", v), nil
	case core.StringType:
		return fmt.Sprintf("%#v", string(v)), nil
	case core.BoolType:
		if v {
			return "true", nil
		}

		return "false", nil
	case core.ErrorType:
		return "", core.Normal(v)
	case core.ListType:
		ts, err := v.ToThunks()

		if err != nil {
			return "", err
		}

		ss := make([]string, 0, len(ts))

		for _, t := range ts {
			s, err := encodeValue(t.Eval())

			if err != nil {
				return "", err
			} else if s == "" {
				continue
			}

			ss = append(ss, s)
		}

		return "[" + strings.Join(ss, ",") + "]", nil
	case core.DictionaryType:
		l, err := core.PApp(core.ToList, core.Normal(v)).EvalList()

		if err != nil {
			return "", core.Normal(err)
		}

		ts, e := l.ToValues()

		if e != nil {
			return "", e
		}

		ss := make([]string, 0, len(ts))

		for _, t := range ts {
			l, err := t.EvalList()

			if err != nil {
				return "", core.Normal(err)
			}

			ts, e := l.ToThunks()

			if e != nil {
				return "", e
			}

			kv := [2]string{}

			for i, t := range ts {
				if i == 0 {
					switch t.Eval().(type) {
					case core.BoolType, core.NilType, core.NumberType:
						s, err := encodeValue(t.Eval())

						if err != nil {
							return "", err
						}

						t = core.NewString(s)
					}
				}

				s, err := encodeValue(t.Eval())

				if err != nil {
					return "", err
				}

				kv[i] = s
			}

			ss = append(ss, kv[0]+":"+kv[1])
		}

		return "{" + strings.Join(ss, ",") + "}", nil
	}

	return "", nil
}
