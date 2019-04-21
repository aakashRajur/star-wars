package types

import (
	"context"
)

type Protocol interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
}
