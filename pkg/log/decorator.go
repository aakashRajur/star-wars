package log

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

func withLocation(logger *logrus.Entry) *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(2); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	} else {
		return logger
	}
}
