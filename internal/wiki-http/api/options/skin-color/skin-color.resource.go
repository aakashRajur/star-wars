package skin_color

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/options/skin-color"
	"github.com/aakashRajur/star-wars/pkg/di/http-resource"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) http_resource.HttpResourceProvider {
	resource := http.NewResource(skin_color.HttpUrl)
	resource.Get(GetSkinColor(storage, logger, tracker))

	return http_resource.HttpResourceProvider{
		Resource: resource,
	}
}
