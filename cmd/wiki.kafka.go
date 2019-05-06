package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/wiki-kafka/module"
	"github.com/aakashRajur/star-wars/pkg/interrupt"
)

func main() {
	defer log.Println(`APPLICATION EXITED`)
	app := fx.New(module.WikiKafkaModule)

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
