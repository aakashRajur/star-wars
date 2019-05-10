package service

type Healthcheck struct {
	URL      string
	Headers  map[string][]string
	Scheme   string
	HttpVerb string
	SkipTLS  bool
	Interval string
	Timeout  string
}
