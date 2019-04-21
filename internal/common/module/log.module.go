package module

import (
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/log"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func GetLogrus() *log.Logger {
	logLevel := env.GetString("LOG_LEVEL")
	logFormatter := env.GetString("LOG_FORMATTER")

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.TraceLevel
	}

	var formatter logrus.Formatter
	switch logFormatter {
	case "json":
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
		}
		break
	case "text":
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
		}
		break
	default:
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339,
		}
	}

	return log.NewInstance(
		level,
		formatter,
	)
}

func GetLogger(logger *log.Logger) types.Logger {
	return logger
}

var LogModule = fx.Provide(
	GetLogrus,
	GetLogger,
)
