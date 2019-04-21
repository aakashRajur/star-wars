package fatal

import (
	"os"
	"testing"

	"github.com/pkg/errors"
)

var osExit = os.Exit

func TestFatalErrorHandler_HandleFatal(t *testing.T) {
	// Save current function and restore at the end:
	defer func() { osExit = os.Exit }()

	var processReturn int
	osExit = func(code int) { processReturn = code }

	handler := FatalErrorHandler{Exit: osExit}
	handler.HandleFatal(errors.New(`silly error`))

	expected := 1

	if processReturn != expected {
		t.Errorf("failed fatal error handling, processReturn: %d, expected: %d", processReturn, expected)
	}
}
