package film

import (
	"context"
	"encoding/json"
	nativeHttp "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/film"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/observable"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

var resourcePatch = film.ResourcePatch

func PatchHttpFilm(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http.HandlerWithMiddleware {
	proxy := http.NewProxy(
		func(url *url.URL, host *url.URL) *url.URL {
			url.Scheme = host.Scheme
			url.Path = strings.Replace(url.Path, httpPrefix, ``, 1)
			url.Host = host.Host
			return url
		},
	)

	requestHandler := func(response http.Response, request *http.Request) {
		defer tracker(time.Now())

		hosts, err := resolver.Resolve(downstreamHttp)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		proxy.UpdateHosts(hosts)
		proxy.HandleRequest(response, request)
	}
	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.Session,
	)
	return http.HandlerWithMiddleware{
		Middlewares:   middlewares,
		HandleRequest: requestHandler,
	}
}

func PatchKafkaFilm(resolver service.Resolver, kafkaInstance *kafka.Kafka, obs *observable.Observable, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics, paramKey string) http.HandlerWithMiddleware {
	requestHandler := func(response http.Response, request *http.Request) {
		defer tracker(time.Now())

		hosts, err := resolver.Resolve(downstreamKafka)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}
		if len(hosts) < 1 {
			response.Error(
				nativeHttp.StatusInternalServerError,
				errors.New(`NO DOWNSTREAM SERVICE AVAILABLE`),
			)
			return
		}

		session, err := middleware.GetSession(request)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		isListening := obs.IsRegistered(session)
		if !isListening {
			response.Error(
				nativeHttp.StatusPreconditionFailed,
				errors.Errorf(`CLIENT %s IS NOT LISTENING ON /async`, session),
			)
			return
		}

		hash, err := util.RandomSHA256()
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		params := request.GetParams()
		id := params[paramKey].(int)

		ctx := request.Context()
		body := ctx.Value(middleware.JSON_BODY).(map[string]interface{})

		payloadListener := make(chan observable.Payload)
		errListener := make(chan error)

		subscription := kafka.Subscription{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  resourcePatch.Type,
			Id:    hash,
			Handler: func(event kafka.Event, instance *kafka.Kafka) {
				var data interface{}

				errMapped := event.Error
				if errMapped != nil {
					data = errMapped
				} else {
					data = event.Id
				}

				marshaled, err := json.Marshal(data)
				if err != nil {
					errListener <- err
					return
				}
				payloadListener <- observable.Payload{
					Key:  event.Id,
					Data: marshaled,
				}
			},
		}

		unsubscriber, err := kafkaInstance.SubscribeOnce(&subscription)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		event := kafka.Event{
			Topic: definedTopics[topics.WikiRequestTopic],
			Type:  resourcePatch.Type,
			Id:    hash,
			Args: map[string]interface{}{
				paramKey: id,
			},
			Data: body,
		}

		timeout, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		rollbackListener, rollback := context.WithCancel(context.Background())
		go func() {
			payload := observable.Payload{
				Key: hash,
			}
			select {
			case <-rollbackListener.Done():
				timeoutCancel()
				err := unsubscriber()
				if err != nil {
					logger.Error(err)
				}
				err = rollbackListener.Err()
				if err != nil && err != context.Canceled {
					logger.Error(err)
				}
				return
			case <-timeout.Done():
				err := unsubscriber()
				if err != nil {
					logger.Error(err)
				}
				err = timeout.Err()
				if err != nil && err != context.DeadlineExceeded {
					payload.Data, err = json.Marshal(err.Error())
					if err != nil {
						logger.Error(err)
						payload.Data = []byte(nativeHttp.StatusText(nativeHttp.StatusInternalServerError))
					}
				}
				break
			case payload = <-payloadListener:
				timeoutCancel()
				break
			case err := <-errListener:
				timeoutCancel()
				payload.Data, err = json.Marshal(err.Error())
				if err != nil {
					logger.Error(err)
					payload.Data = []byte(nativeHttp.StatusText(nativeHttp.StatusInternalServerError))
				}
			}

			err = obs.SendData(session, payload)
			if err != nil {
				logger.Error(err)
			}
		}()

		err = kafkaInstance.Emit(event)
		if err != nil {
			rollback()
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		err = response.WriteJSON(hash, nil)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.Session,
		middleware.ValidateArgs(
			film.ArgValidation,
			film.ArgNormalization,
		),
		middleware.JsonBodyParser,
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
