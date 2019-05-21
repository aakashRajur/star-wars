package validations

import (
	"testing"

	"github.com/pkg/errors"
)

func TestRequired(t *testing.T) {
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
		`DOES_NOT_EXIST`: {
			Args: arg{
				Key:    `key1`,
				Value:  nil,
				Exists: false,
			},
			Expected: errors.Errorf(RequiredError, `key1`),
		},
		`NIL`: {
			Args: arg{
				Key:    `key2`,
				Value:  nil,
				Exists: true,
			},
			Expected: errors.Errorf(RequiredError, `key2`),
		},
		`VALID`: {
			Args: arg{
				Key:    `key3`,
				Value:  `hello_world`,
				Exists: true,
			},
			Expected: nil,
		},
	}
	validator := ValidateRequired()

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
			t.Errorf("ValidateRequired() = %+v, want %+v", e1, e2)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
