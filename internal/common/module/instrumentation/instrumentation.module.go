package instrumentation

import (
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

		logger.InfoFields(
			map[string]interface{}{
				`function`: funcObj.Name(),
				`elapsed`:  elapsed,
			},
		)
	}
}

var Module = fx.Provide(
	GetExecutionTracer,
)
