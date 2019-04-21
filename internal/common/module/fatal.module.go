package module

import (
	"os"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/fatal"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetFatalErrorHandler() types.FatalHandler {
	return fatal.FatalErrorHandler{Exit: os.Exit}
}

var FatalHandlerModule = fx.Provide(GetFatalErrorHandler)
