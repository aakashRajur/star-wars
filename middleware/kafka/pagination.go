package kafka

import (
	"context"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Pagination() kafka.Middleware {
	return func(next kafka.HandleEvent) kafka.HandleEvent {
		return func(event kafka.Event, instance *kafka.Kafka) {
			pagination := types.Pagination{
				PaginationId: 0,
				Limit:        10,
			}

			args := event.Args
			if args == nil {
				newCtx := context.WithValue(event.Ctx, types.PAGINATION, pagination)
				event.Ctx = newCtx
				next(event, instance)
				return
			}

			paginationId, ok := args[types.QueryPaginationId]
			if ok {
				parsed, ok := paginationId.(int64)
				if ok {
					pagination.PaginationId = parsed
				}
			}

			limit, ok := args[types.QueryLimit]
			if ok {
				parsed, ok := limit.(int64)
				if ok {
					pagination.Limit = parsed
				}
			}

			newCtx := context.WithValue(event.Ctx, types.PAGINATION, pagination)
			event.Ctx = newCtx
			next(event, instance)
		}
	}
}
