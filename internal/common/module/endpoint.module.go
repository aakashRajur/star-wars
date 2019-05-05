package module

import (
	"fmt"
	"github.com/aakashRajur/star-wars/pkg/env"
	"go.uber.org/fx"
)

func GetEndpoint() string {
	return fmt.Sprintf(`%s:%s`, env.GetString(`CONTAINER_HOST_NAME`), env.GetString(`CONTAINER_PORT`))
}

var EndpointModule = fx.Provide(GetEndpoint)
