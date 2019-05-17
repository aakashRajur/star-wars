package main

import (
	"context"
	"github.com/aakashRajur/star-wars/internal/wiki-kafka/module/app"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/interrupt"
)

func main() {
	defer log.Println(`APPLICATION EXITED`)
	application := fx.New(app.WikiKafkaModule)

	wait := interrupt.NotifyOnInterrupt(
		application.Stop,
		60*time.Second,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := application.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-wait
}
