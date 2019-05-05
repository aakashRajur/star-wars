package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/wiki/module/kafka"
	"github.com/aakashRajur/star-wars/pkg/interrupt"
)

func main() {
	defer log.Println(`Exiting...`)
	app := fx.New(kafka.WikiKafkaModule)

	wait := interrupt.NotifyOnInterrupt(
		app.Stop,
		30*time.Second,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-wait
}
