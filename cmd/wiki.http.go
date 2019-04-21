package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/wiki/module/http"
	"github.com/aakashRajur/star-wars/pkg/interrupt"
)

func main() {
	defer log.Println(`Exiting...`)
	app := fx.New(http.WikiHttpModule)

	wait := interrupt.NotifyOnInterrupt(
		app.Stop,
		30*time.Second,
		os.Interrupt,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-wait
}
