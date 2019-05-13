package service

import "context"

type HandleEvent func(event interface{}, ctx context.Context)

type Subscription struct {
	Key     string
	Handler HandleEvent
}

type Observer interface {
	Observe(subscription Subscription) error
}
