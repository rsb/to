package to

import (
	"encoding/json"
	"github.com/rsb/failure"
	"golang.org/x/exp/constraints"
	"reflect"
	"strconv"
	"time"
)

type IntData[T constraints.Signed] struct {
	item     *T
	sValue   string
	typeName string
}

func NewIntData[T constraints.Signed](v *T) IntData[T] {
	t := reflect.TypeOf(v)
	n := IntData[T]{
		item:     v,
		typeName: t.Name(),
	}

	return n
}

func (d *IntData[T]) Item() *T {
	return d.item
}

func (d *IntData[T]) Set(v string) error {
	i, err := Int[T](v)
	if err != nil {
		return failure.Wrap(err, "Int[%v] failed", d.typeName)
	}

	d.item = &i
	return nil
}

func (d *IntData[T]) Type() string {
	return d.typeName
}

func (d *IntData[T]) String() string {
	return String(d.item)
}

func Int[T constraints.Signed](i any) (T, error) {
	i = indirect(i)

	v, ok := integer(i)
	if ok {
		return T(v), nil
	}

	switch s := i.(type) {
	case int8:
		return T(s), nil
	case int16:
		return T(s), nil
	case int32:
		return T(s), nil
	case int64:
		return T(s), nil
	case uint:
		return T(s), nil
	case uint8:
		return T(s), nil
	case uint16:
		return T(s), nil
	case uint32:
		return T(s), nil
	case uint64:
		return T(s), nil
	case float32:
		return T(s), nil
	case float64:
		return T(s), nil
	case string:
		v, err := strconv.ParseInt(s, 0, 0)
		if err != nil {
			return 0, failure.ToInvalidParam(err, "unable to cast %#v of type %T to int64", i, i)
		}
		return T(v), nil
	case json.Number:
		v, err := Int[T](string(s))
		if err != nil {
			return 0, failure.ToInvalidParam(err, "Int failed for json.Number (%v)", i)
		}
		return v, nil
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, failure.InvalidParam("unable to cast %#v of type %T to int", i, i)
	}
}

// integer returns the int value of v if v or v's underlying type
// is an int.
// Note that this will return false for int64 etc. types.
func integer(v any) (int, bool) {
	switch v := v.(type) {
	case int:
		return v, true
	case time.Weekday:
		return int(v), true
	case time.Month:
		return int(v), true
	default:
		return 0, false
	}
}
