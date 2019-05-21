package validations

import (
	"fmt"
	"testing"

	"github.com/juju/errors"
)

func TestValidateInterval(t *testing.T) {
	type arg struct {
		Key    string
		Value  interface{}
		Exists bool
	}
	type testCase struct {
		Args     arg
		Expected error
	}

	testCases := map[string]testCase{
		`INCORRECT_FIELD_VALUE_TYPE`: {
			Args: arg{
				Key: `key1`,
				Value: map[string]interface{}{
					`year`:  `1`,
					`month`: 10,
				},
				Exists: true,
			},
			Expected: errors.Errorf(IntegerError, `1`, `year`),
		},
		`EXTRA_KEYS`: {
			Args: arg{
				Key: `key1`,
				Value: map[string]interface{}{
					`year`:  1,
					`month`: 10,
					`hello`: `world`,
				},
				Exists: true,
			},
			Expected: errors.Errorf(
				`{"hello":"%s"}`,
				fmt.Sprintf(IntervalError, `hello`),
			),
		},
		`VALID_INTERVAL`: {
			Args: arg{
				Key: `key2`,
				Value: map[string]interface{}{
					`year`:  1,
					`month`: 2,
					`day`:   3,
					`hour`:  4,
					`min`:   5,
					`sec`:   6,
				},
				Exists: true,
			},
			Expected: nil,
		},
	}

	validator := ValidateInterval()

	for name, test := range testCases {
		args := test.Args
		got := validator(args.Key, args.Value, args.Exists)
		expected := test.Expected

		e1 := `nil`
		e2 := `nil`

		if got != nil {
			e1 = got.Error()
		}
		if expected != nil {
			e2 = expected.Error()
		}

		if e1 != e2 {
			t.Errorf("ValidateInterval() = %+v, want %+v", e1, e2)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
