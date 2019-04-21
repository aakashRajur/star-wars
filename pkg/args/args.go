package args

import (
	"flag"
	"os"
)

// will load app env (dev|prod)
func LoadArgs() {
	envPtr := flag.String("env", "dev", "env type")

	flag.Parse()
	_ = os.Setenv("env", *envPtr)
}
