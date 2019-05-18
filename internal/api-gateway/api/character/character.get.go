package character

import (
	"context"
	"encoding/json"
	"fmt"
	nativeHttp "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/internal/topics"
	"github.com/aakashRajur/star-wars/internal/wiki/api/character"
	middleware "github.com/aakashRajur/star-wars/middleware/http"
	"github.com/aakashRajur/star-wars/pkg/http"
	"github.com/aakashRajur/star-wars/pkg/kafka"
	"github.com/aakashRajur/star-wars/pkg/observable"
	"github.com/aakashRajur/star-wars/pkg/service"
	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

var resourceGet = character.ResourceGet

func GetHttpCharacter(resolver service.Resolver, logger types.Logger, tracker types.TimeTracker) http.HandlerWithMiddleware {
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

func GetKafkaCharacter(resolver service.Resolver, kafkaInstance *kafka.Kafka, observable *observable.Observable, logger types.Logger, tracker types.TimeTracker, definedTopics kafka.DefinedTopics, paramKey string) http.HandlerWithMiddleware {
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

		ctx := request.Context()
		sessionValue := ctx.Value(middleware.SESSION_COOKIE)
		if sessionValue == nil {
			response.Error(
				nativeHttp.StatusInternalServerError,
				errors.New(`NO SESSION SET`),
			)
			return
		}
		session, ok := sessionValue.(string)
		if !ok {
			response.Error(
				nativeHttp.StatusInternalServerError,
				errors.New(`CORRUPT SESSION`),
			)
			return
		}

		isListening := observable.IsRegistered(session)
		if !isListening {
			response.Error(
				nativeHttp.StatusPreconditionFailed,
				errors.Errorf(`CLIENT %s IS NOT LISTENING ON /async`, session),
			)
			return
		}

		hash, err := util.SHA256()
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		params := request.GetParams()
		id := params[paramKey].(int)

		dataListener := make(chan []byte)
		errListener := make(chan error)

		subscription := kafka.Subscription{
			Topic: definedTopics[topics.WikiResponseTopic],
			Type:  resourceGet.Type,
			Id:    hash,
			Handler: func(event kafka.Event, instance *kafka.Kafka) {
				var data interface{}

				errMapped := event.Error
				if errMapped != nil {
					data = errMapped
				} else {
					data = event.Data
				}

				marshaled, err := json.Marshal(data)
				if err != nil {
					errListener <- err
					return
				}
				dataListener <- marshaled
			},
		}

		unsubscriber, err := kafkaInstance.SubscribeOnce(&subscription)
		if err != nil {
			response.Error(nativeHttp.StatusInternalServerError, err)
			return
		}

		event := kafka.Event{
			Topic: definedTopics[topics.WikiRequestTopic],
			Type:  resourceGet.Type,
			Id:    hash,
			Args: map[string]interface{}{
				paramKey: id,
			},
			Data: nil,
		}

		timeout, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
		rollbackListener, rollback := context.WithCancel(context.Background())
		go func() {
			var data []byte
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
					data, err = json.Marshal(err.Error())
					if err != nil {
						logger.Error(err)
						data = []byte(nativeHttp.StatusText(nativeHttp.StatusInternalServerError))
					}
				}
				break
			case data = <-dataListener:
				timeoutCancel()
				break
			case err := <-errListener:
				timeoutCancel()
				data, err = json.Marshal(err.Error())
				if err != nil {
					logger.Error(err)
					data = []byte(nativeHttp.StatusText(nativeHttp.StatusInternalServerError))
				}
			}

			err = observable.SendData(session, []byte(fmt.Sprintf(`id: %s`, session)))
			if err != nil {
				logger.Error(err)
				return
			}
			err = observable.SendData(session, data)
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

		err = response.WriteJSON(id, nil)
		if err != nil {
			logger.Error(err)
		}
	}

	middlewares := http.ChainMiddlewares(
		middleware.Logger(logger),
		middleware.Session,
		middleware.ValidateArgs(
			character.ArgValidation,
			character.ArgNormalization,
		),
	)

	return http.HandlerWithMiddleware{
		HandleRequest: requestHandler,
		Middlewares:   middlewares,
	}
}
