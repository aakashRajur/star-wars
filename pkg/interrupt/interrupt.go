package interrupt

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func NotifyOnInterrupt(onClose func(context.Context) error, timeout time.Duration, signals ...os.Signal) chan bool {
	wait := make(chan bool)

	go func() {
		sigint := make(chan os.Signal, 1)
		for i, iL := 0, len(signals); i < iL; i++ {
			signal.Notify(sigint, signals[i])
		}
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		if err := onClose(ctx); err != nil {
			log.Fatal(err)
		}
		cancel()
		close(wait)
	}()

	return wait
}
