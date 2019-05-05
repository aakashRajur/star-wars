package interrupt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func NotifyOnInterrupt(onClose func(context.Context) error, timeout time.Duration) chan bool {
	wait := make(chan bool)

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt)
		<-interrupt

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		if err := onClose(ctx); err != nil {
			log.Fatal(err)
		}
		cancel()
		close(wait)
	}()

	return wait
}
