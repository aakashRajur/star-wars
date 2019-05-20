package redis

import (
	"github.com/go-redis/redis"

	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

func NewLogger(logger types.Logger) func(func(redis.Cmder) error) func(redis.Cmder) error {
	return func(oldProcess func(redis.Cmder) error) func(redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			err := oldProcess(cmd)
			args := cmd.Args()

			logFields := map[string]interface{}{`cmd`: args[0]}
			safe, err := util.MapStringFromInterfaces(args[1:])
			if err != nil {
				logFields[`args`] = args[1:]
			} else {
				logFields[`args`] = safe
			}

			if err != nil {
				logger.ErrorFields(err, logFields)
			} else {
				logger.InfoFields(logFields)
			}

			return err
		}
	}
}
