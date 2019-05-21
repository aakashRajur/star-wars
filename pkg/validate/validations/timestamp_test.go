package validations

import (
	"testing"

	"github.com/juju/errors"
)

func TestValidateTimestamp(t *testing.T) {
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
		`VALID_TIMESTAMP`: {
			Args: arg{
				Key:    `key1`,
				Value:  `2019-05-21T16:17:51.062Z`,
				Exists: true,
			},
			Expected: nil,
		},
		`INVALID_TIMESTAMP`: {
			Args: arg{
				Key:    `key2`,
				Value:  `hello_world`,
				Exists: true,
			},
			Expected: errors.New(`parsing time "hello_world" as "2006-01-02T15:04:05Z0700": cannot parse "hello_world" as "2006"`),
		},
	}

	validator := ValidateTimestamp()

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
			t.Errorf("ValidateTimestamp() = %+v, want %+v", e1, e2)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
