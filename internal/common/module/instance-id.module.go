package module

import (
	"errors"
	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/types"
	"go.uber.org/fx"
)

func GetInstanceId(handler types.FatalHandler) types.InstanceId {
	instanceId := env.GetString(`INSTANCE_ID`)
	if instanceId == `` {
		handler.HandleFatal(
			errors.New(`INSTANCE ID CANNOT BE EMPTY`),
		)
		return ``
	}

	return types.InstanceId(instanceId)
}

var InstanceIdModule = fx.Provide(GetInstanceId)
