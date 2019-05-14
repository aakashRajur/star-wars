package service

import "context"

type HandleEvent func(event interface{}, ctx context.Context)

type Subscription struct {
	Key     string
	Handler HandleEvent
}
