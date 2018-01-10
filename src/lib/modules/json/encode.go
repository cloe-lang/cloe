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
	func(ts ...core.Value) core.Value {
		s, err := encodeValue(ts[0].Eval())

		if err != nil {
			return err
		} else if s == "" {
			return jsonEncodeError
		}

		return core.NewString(s)
	})

func encodeValue(v core.Value) (string, core.Value) {
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
		return "", v
	case core.ListType:
		ss := []string{}

		for !v.Empty() {
			s, err := encodeValue(v.First().Eval())

			if err != nil {
				return "", err
			} else if s != "" {
				ss = append(ss, s)
			}

			v, err = v.Rest().EvalList()

			if err != nil {
				return "", err
			}
		}

		return "[" + strings.Join(ss, ",") + "]", nil
	case core.DictionaryType:
		l, err := core.PApp(core.ToList, v).EvalList()

		if err != nil {
			return "", err
		}

		ss := []string{}

		for !l.Empty() {
			ll, err := l.First().EvalList()

			if err != nil {
				return "", err
			}

			kk := ll.First().Eval()

			switch kk.(type) {
			case core.BoolType, core.NilType, core.NumberType:
				s, err := encodeValue(kk)

				if err != nil {
					return "", err
				}

				kk = core.StringType(s)
			}

			k, err := encodeValue(kk)

			if err != nil {
				return "", err
			}

			ll, err = ll.Rest().EvalList()

			if err != nil {
				return "", err
			}

			v, err := encodeValue(ll.First().Eval())

			if err != nil {
				return "", err
			}

			ss = append(ss, k+":"+v)

			l, err = l.Rest().EvalList()

			if err != nil {
				return "", err
			}
		}

		return "{" + strings.Join(ss, ",") + "}", nil
	}

	return "", nil
}
