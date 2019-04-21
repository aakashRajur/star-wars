package fatal

import (
	"fmt"

	"github.com/juju/errors"
)

//noinspection GoNameStartsWithPackageName
type FatalErrorHandler struct {
	Exit func(int)
}

func (errorHandler FatalErrorHandler) HandleFatal(err error) {
	fmt.Println(errors.Trace(err))
	errorHandler.Exit(1)
}
