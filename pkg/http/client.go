package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const (
	VerbGet    = `GET`
	VerbPatch  = `PATCH`
	VerbPut    = `PUT`
	VerbDelete = `DELETE`
)

func withParams(u *url.URL, params map[string]interface{}) string {
	path := u.Path
	for key, value := range params {
		path = strings.Replace(
			path,
			fmt.Sprintf(`:%s`, key),
			fmt.Sprintf(`%v`, value),
			-1,
		)
	}

	formatted := &url.URL{
		Scheme:     u.Scheme,
		Opaque:     u.Opaque,
		User:       u.User,
		Host:       u.Host,
		Path:       path,
		RawPath:    u.RawPath,
		ForceQuery: u.ForceQuery,
		RawQuery:   u.RawQuery,
		Fragment:   u.Fragment,
	}

	return formatted.String()
}

func NewRequest(config RequestConfig) (*http.Response, error) {
	var body io.Reader
	if config.Body != nil {
		marshaled, err := json.Marshal(config.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(marshaled)
	}

	client := http.Client{Timeout: config.Timeout}

	req, err := http.NewRequest(
		config.Verb,
		withParams(&config.Url, config.Params),
		body,
	)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func NewCanceleableRequest(config RequestConfig, ctx context.Context) (*http.Response, error) {
	safeCtx := ctx
	if safeCtx != nil {
		safeCtx = context.Background()
	}
	responseChan := make(chan *http.Response)
	errorChan := make(chan error)

	go func(config RequestConfig, rChan chan<- *http.Response, eChan chan<- error) {
		response, err := NewRequest(config)
		if err != nil {
			eChan <- err
			return
		}
		rChan <- response
	}(config, responseChan, errorChan)

	select {
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return nil, err
		}
	case err := <-errorChan:
		return nil, err
	case response := <-responseChan:
		return response, nil
	}
	return nil, errors.New(`INVALID HTTP REQUEST STATE`)
}

func DataFromResponse(response *http.Response) ([]byte, error) {
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	if len(data) < 1 {
		return []byte(`N/A`), nil
	}

	return data, nil
}

func TextFromResponse(response *http.Response) (string, error) {
	status := response.StatusCode
	if status < 200 || status > 299 {
		body, err := DataFromResponse(response)
		if err != nil {
			return ``, errors.Errorf(
				`REQUEST STATUS %d, FAILED TO PARSE ERROR BODY: %s`,
				status,
				err.Error(),
			)
		} else {
			return ``, errors.Errorf(
				`REQUEST STATUS %d, ERROR BODY: %s`,
				status,
				string(body),
			)
		}
	}
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ``, err
	}

	return string(data), nil
}

func JsonObjectFromResponse(response *http.Response) (map[string]interface{}, error) {
	status := response.StatusCode
	if status < 200 || status > 299 {
		body, err := DataFromResponse(response)
		if err != nil {
			return nil, errors.Errorf(
				`REQUEST STATUS %d, FAILED TO PARSE ERROR BODY: %s`,
				status,
				err.Error(),
			)
		} else {
			return nil, errors.Errorf(
				`REQUEST STATUS %d, ERROR BODY: %s`,
				status,
				body,
			)
		}
	}
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsed := make(map[string]interface{})
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func JsonArrayFromResponse(response *http.Response) ([]map[string]interface{}, error) {
	status := response.StatusCode
	if status < 200 || status > 299 {
		body, err := DataFromResponse(response)
		if err != nil {
			return nil, errors.Errorf(
				`REQUEST STATUS %d, FAILED TO PARSE ERROR BODY: %s`,
				status,
				err.Error(),
			)
		} else {
			return nil, errors.Errorf(
				`REQUEST STATUS %d, ERROR BODY: %s`,
				status,
				string(body),
			)
		}
	}
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	parsed := make([]map[string]interface{}, 0)
	err = json.Unmarshal(data, &parsed)
	if err != nil {
		return nil, err
	}

	return parsed, nil
}
