package json

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cloe-lang/cloe/src/lib/core"
)

var jsonEncodeError = jsonError(errors.New("cannot be encoded as JSON"))

var encode = core.NewLazyFunction(
	core.NewSignature([]string{"decoded"}, "", nil, ""),
	func(vs ...core.Value) core.Value {
		s, err := encodeValue(vs[0])

		if err != nil {
			return err
		} else if s == "" {
			return jsonEncodeError
		}

		return core.NewString(s)
	})

func encodeValue(v core.Value) (string, core.Value) {
	switch v := core.EvalPure(v).(type) {
	case core.NilType:
		return "null", nil
	case *core.NumberType:
		return fmt.Sprintf("%v", *v), nil
	case core.StringType:
		return fmt.Sprintf("%#v", string(v)), nil
	case *core.BooleanType:
		if *v {
			return "true", nil
		}

		return "false", nil
	case *core.ErrorType:
		return "", v
	case *core.ListType:
		ss := []string{}

		for !v.Empty() {
			s, err := encodeValue(v.First())

			if err != nil {
				return "", err
			} else if s != "" {
				ss = append(ss, s)
			}

			v, err = core.EvalList(v.Rest())

			if err != nil {
				return "", err
			}
		}

		return "[" + strings.Join(ss, ",") + "]", nil
	case *core.DictionaryType:
		l, err := core.EvalList(core.PApp(core.ToList, v))

		if err != nil {
			return "", err
		}

		ss := []string{}

		for !l.Empty() {
			ll, err := core.EvalList(l.First())

			if err != nil {
				return "", err
			}

			kk := core.EvalPure(ll.First())

			switch kk.(type) {
			case *core.BooleanType, core.NilType, *core.NumberType:
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

			ll, err = core.EvalList(ll.Rest())

			if err != nil {
				return "", err
			}

			v, err := encodeValue(ll.First())

			if err != nil {
				return "", err
			}

			ss = append(ss, k+":"+v)

			l, err = core.EvalList(l.Rest())

			if err != nil {
				return "", err
			}
		}

		return "{" + strings.Join(ss, ",") + "}", nil
	}

	return "", nil
}
