package http

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ContentType         = "Content-Type"
	ContentTypeJSON     = "application/json"
	ContentTypeText     = "application/text"
	ContentLength       = `Content-Length`
	ContentEncoding     = `Content-Encoding`
	ContentEncodingGzip = `gzip`
)

type Response struct {
	http.ResponseWriter
	compress bool
}

func (response Response) Header() http.Header {
	return response.ResponseWriter.Header()
}

func (response Response) WriteHeader(status int) {
	response.ResponseWriter.WriteHeader(status)
}

func (response Response) Write(data []byte) (int, error) {
	var writer io.Writer = response.ResponseWriter
	if response.compress {
		header := response.ResponseWriter.Header()
		header.Del(ContentLength)
		header.Set(ContentEncoding, ContentEncodingGzip)

		gz := gzipPool.Get().(*gzip.Writer)
		defer gzipPool.Put(gz)

		gz.Reset(writer)
		defer gz.Close()
		writer = gz
	}
	return writer.Write(data)
}

func (response Response) SetStatus(status int) Response {
	response.WriteHeader(status)
	return response
}

func (response Response) GetStatus() string {
	return response.GetHeader("status")
}

func (response Response) SetHeader(key string, value string) Response {
	response.ResponseWriter.Header().Add(key, value)
	return response
}

func (response Response) GetHeader(key string) string {
	return response.Header().Get(key)
}

func (response Response) DelHeader(key string) Response {
	response.Header().Del(key)
	return response
}

func (response Response) WriteText(text string) (int, error) {
	return fmt.Fprintf(response, text)
}

func (response Response) WriteJSON(data interface{}, headers map[string]string) error {
	response.SetHeader(ContentType, ContentTypeJSON)
	if headers != nil {
		for key, value := range headers {
			response.SetHeader(key, value)
		}
	}
	return json.NewEncoder(response).Encode(data)
}

func (response Response) WriteJSONIndent(data interface{}, headers map[string]string) error {
	response.SetHeader(ContentType, ContentTypeJSON)
	if headers != nil {
		for key, value := range headers {
			response.SetHeader(key, value)
		}
	}
	encoder := json.NewEncoder(response)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (response Response) Error(status int, err error) {
	http.Error(
		response.ResponseWriter,
		err.Error(),
		status,
	)
}

func (response Response) ErrorString(status int, contentType string, err string) {
	response.SetHeader(ContentType, contentType)
	http.Error(
		response.ResponseWriter,
		err,
		status,
	)
}
