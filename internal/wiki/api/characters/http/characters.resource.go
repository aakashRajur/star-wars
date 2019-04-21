package characters

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/characters"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.ResourceProvider {
	resource := http.NewResource(characters.HttpURI)
	resource.Get(GetCharacters(storage, logger, tracker, characters.CacheKey))

	return di.ResourceProvider{
		Resource: resource,
	}
}
