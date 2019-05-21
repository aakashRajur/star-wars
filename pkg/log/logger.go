package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Logger.Infoln(args...)
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Traceln(args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Debugln(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	withLocation(logrus.NewEntry(logger.Logger)).Warningln(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	withLocation(logrus.NewEntry(logger.Logger)).Errorln(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	withLocation(logrus.NewEntry(logger.Logger)).Fatalln(args...)
}

func (logger *Logger) InfoFields(fields map[string]interface{}, args ...interface{}) {
	logger.WithFields(fields).Infoln(args...)
}

func (logger *Logger) TraceFields(fields map[string]interface{}, args ...interface{}) {
	logger.WithFields(fields).Traceln(args...)
}

func (logger *Logger) DebugFields(fields map[string]interface{}, args ...interface{}) {
	logger.WithFields(fields).Debugln(args...)
}

func (logger *Logger) WarnFields(fields map[string]interface{}, args ...interface{}) {
	withLocation(logger.WithFields(fields)).Warnln(args...)
}

func (logger *Logger) ErrorFields(err error, fields map[string]interface{}, args ...interface{}) {
	withLocation(logger.WithError(err).WithFields(fields)).Errorln(args...)
}

func (logger *Logger) FatalFields(fields map[string]interface{}, args ...interface{}) {
	withLocation(logger.WithFields(fields)).Fatalln(args...)
}

func (logger *Logger) Out() io.Writer {
	return logger.Logger.Out
}

func NewInstance(level logrus.Level, formatter logrus.Formatter) *Logger {
	enhanced := logrus.New()
	enhanced.SetLevel(level)
	enhanced.SetFormatter(formatter)

	logger := Logger{enhanced}
	return &logger
}
