package types

import (
	"io"
)

type Logger interface {
	Info(args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	InfoFields(fields map[string]interface{}, args ...interface{})
	TraceFields(fields map[string]interface{}, args ...interface{})
	DebugFields(fields map[string]interface{}, args ...interface{})
	WarnFields(fields map[string]interface{}, args ...interface{})
	ErrorFields(err error, fields map[string]interface{}, args ...interface{})
	FatalFields(fields map[string]interface{}, args ...interface{})
	Out() io.Writer
}
