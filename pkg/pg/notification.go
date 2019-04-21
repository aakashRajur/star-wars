package pg

import (
	"time"

	"github.com/jackc/pgx"
)

type Notification struct {
	*pgx.Notification
	timestamp time.Time
}

func (notification *Notification) GetChannel() string {
	return notification.Channel
}

func (notification *Notification) GetPayload() string {
	return notification.Payload
}

func (notification *Notification) GetTimestamp() time.Time {
	return notification.timestamp
}

func NewNotification(original *pgx.Notification) *Notification {
	return &Notification{
		Notification: original,
		timestamp:    time.Now(),
	}
}
