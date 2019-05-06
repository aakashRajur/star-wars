package module

import (
	"fmt"

	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/types"

	"go.uber.org/fx"
)

func GetEndpoint() types.Endpoint {
	endpoint := fmt.Sprintf(`%s:%s`, env.GetString(`CONTAINER_HOST_NAME`), env.GetString(`CONTAINER_PORT`))
	return types.Endpoint(endpoint)
}

var EndpointModule = fx.Provide(GetEndpoint)
