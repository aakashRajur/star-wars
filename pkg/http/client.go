package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	VerbGet   = `GET`
	VerbPatch = `PATCH`
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

func JsonObject(response *http.Response) (map[string]interface{}, error) {
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

func JsonArray(response *http.Response) ([]map[string]interface{}, error) {
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
