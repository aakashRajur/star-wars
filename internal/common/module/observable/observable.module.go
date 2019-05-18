package observable

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/pkg/observable"
)

func GetObservable() *observable.Observable {
	return observable.NewInstance()
}

var Module = fx.Provide(GetObservable)
