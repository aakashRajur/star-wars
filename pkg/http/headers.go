package http

const (
	ContentType           = "Content-Type"
	ContentTypeJSON       = "application/json"
	ContentTypeText       = "application/text"
	ContentTypeTextStream = `text/event-stream`
	ContentLength         = `Content-Length`
	ContentEncoding       = `Content-Encoding`
	CacheControl          = `Cache-Control`
	CacheControlNoCache   = `no-cache`
	Connection            = `Connection`
	ConnectionKeepAlive   = `keep-alive`
	ContentEncodingGzip   = `gzip`
	XForwardedFor         = `X-Forwarded-For`
)
