package pg

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

type pgLogger struct {
	types.Logger
}

func (logger pgLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	native := logger.Logger
	switch level {
	case pgx.LogLevelTrace:
		native.TraceFields(data, msg)
	case pgx.LogLevelDebug:
		native.DebugFields(data, msg)
	case pgx.LogLevelInfo:
		native.InfoFields(data, msg)
	case pgx.LogLevelWarn:
		native.WarnFields(data, msg)
	case pgx.LogLevelError:
		native.ErrorFields(errors.New(msg), data)
	}
}

func NewPgLogger(logger types.Logger) pgx.Logger {
	fromInterface := pgLogger{logger}
	return fromInterface
}
