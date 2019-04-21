package util

import (
	"reflect"
	"strconv"

	"github.com/juju/errors"
)

func MapStringFromInterfaces(args []interface{}) ([]string, error) {
	compiled := make([]string, 0)
	for _, each := range args {
		switch t := each.(type) {
		case float64:
			{
				compiled = append(compiled, strconv.FormatFloat(t, 'f', 64, -1))
				break
			}
		case float32:
			{
				compiled = append(compiled, strconv.FormatFloat(float64(t), 'f', 64, -1))
				break
			}
		case int64:
			{
				compiled = append(compiled, strconv.FormatInt(t, 10))
				break
			}
		case int32:
		case int16:
		case int8:
		case int:
			{
				compiled = append(compiled, strconv.FormatInt(int64(t), 10))
				break
			}
		case bool:
			{
				compiled = append(compiled, strconv.FormatBool(t))
				break
			}
		case string:
			{
				compiled = append(compiled, t)
				break
			}
		case []byte:
			{
				compiled = append(compiled, string(t[:]))
			}
		default:
			return nil, errors.New(`type not supported`)
		}
	}
	return compiled, nil
}

func GetType(value interface{}) string {
	return reflect.TypeOf(value).Kind().String()
}
