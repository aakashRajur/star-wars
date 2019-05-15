package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	VerbGet   = `GET`
	VerbPatch = `PATCH`
	VerbPut   = `PUT`
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

func NewRequest(config RequestConfig) (*http.Request, error) {
	var body io.Reader
	if config.Body != nil {
		marshaled, err := json.Marshal(config.Body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(marshaled)
	}

	return http.NewRequest(
		config.Verb,
		withParams(&config.Url, config.Params),
		body,
	)
}
