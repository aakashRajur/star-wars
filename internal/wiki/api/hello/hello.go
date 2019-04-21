package hello

import (
	"fmt"
	"time"
)

func SayHello() string {
	return fmt.Sprintf("hello world!, %s", time.Now().UTC())
}
