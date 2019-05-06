package module

import (
	"go.uber.org/fx"

	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/pkg/env"
	"github.com/aakashRajur/star-wars/pkg/kafka"
)

func GetDefinedTopics() kafka.DefinedTopics {
	var definedTopics kafka.DefinedTopics = make(map[string]string)

	definedTopics[topics.WikiRequestTopic] = env.GetString(topics.WikiRequestTopic)
	definedTopics[topics.WikiResponseTopic] = env.GetString(topics.WikiResponseTopic)
	definedTopics[topics.AuthRequestTopic] = env.GetString(topics.AuthRequestTopic)
	definedTopics[topics.AuthResponseTopic] = env.GetString(topics.AuthResponseTopic)
	definedTopics[topics.SearchRequestTopic] = env.GetString(topics.SearchRequestTopic)
	definedTopics[topics.SearchResponseTopic] = env.GetString(topics.SearchResponseTopic)

	return definedTopics
}

var TopicsModule = fx.Provide(GetDefinedTopics)
