package character

import (
	"github.com/aakashRajur/star-wars/internal/wiki/api/character"
	"github.com/aakashRajur/star-wars/internal/wiki/api/characters"
	"github.com/aakashRajur/star-wars/pkg/di"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/types"
)

func Resource(storage types.Storage, logger types.Logger, tracker types.TimeTracker) di.ResourceProvider {
	resource := http.NewResource(character.HttpURL)
	resource.Get(GetCharacter(storage, logger, tracker, characters.CacheKey, character.ParamCharacterId))
	resource.Patch(PatchCharacter(storage, logger, tracker, character.ParamCharacterId))

	return di.ResourceProvider{Resource: resource}
}
