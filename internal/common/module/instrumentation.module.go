package module

import (
	"regexp"
	"runtime"
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetExecutionTracer(logger types.Logger) types.TimeTracker {
	return func(start time.Time) {
		elapsed := time.Since(start)

		// Skip this function, and fetch the PC and file for its parent.
		pc, _, _, _ := runtime.Caller(1)

		// Retrieve a function object this functions parent.
		funcObj := runtime.FuncForPC(pc)

		// Regex to extract just the function name (and not the common path).
		runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
		name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

		logger.InfoFields(
			map[string]interface{}{
				`function`: name,
				`elapsed`:  elapsed,
			},
		)
	}
}

var InstrumentationModule = fx.Provide(
	GetExecutionTracer,
)
