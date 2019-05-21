package args

import (
	"flag"
	"os"
	"testing"
)

type testCase struct {
	Args     []string
	Expected string
}

var testCases = map[string]testCase{
	`DEFAULT`: {
		Expected: `dev`,
		Args:     []string{"cmd"},
	},
	`DEV`: {
		Expected: `dev`,
		Args:     []string{"cmd", "-env=dev"},
	},
	`PROD`: {
		Expected: `prod`,
		Args:     []string{"cmd", "-env=prod"},
	},
}

func TestLoadArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for name, tc := range testCases {
		flag.CommandLine = flag.NewFlagSet(tc.Args[0], flag.ContinueOnError)
		flag.CommandLine.Usage = flag.Usage
		os.Args = tc.Args

		LoadArgs()
		value := os.Getenv(`env`)
		if value != tc.Expected {
			t.Errorf(`%s failed, got: %s, expected: %s`, name, value, tc.Expected)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
