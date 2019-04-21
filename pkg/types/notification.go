package types

import (
	"time"
)

type Notification interface {
	GetChannel() string
	GetPayload() string
	GetTimestamp() time.Time
}
