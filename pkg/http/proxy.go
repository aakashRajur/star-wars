package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/juju/errors"
)

func updateHeaders(header http.Header, additional map[string][]string) http.Header {
	for key, values := range additional {
		for _, value := range values {
			header.Add(key, value)
		}
	}

	return header
}

type Proxy struct {
	downstreamUrl func(url *url.URL, host *url.URL) *url.URL
	hosts         map[string]int
}

func (proxy *Proxy) UpdateHosts(hosts []string) {
	existing := proxy.hosts
	tracker := make(map[string]bool)

	for _, each := range hosts {
		tracker[each] = true
		_, ok := existing[each]
		if !ok {
			existing[each] = 0
		}
	}
	purge := make([]string, 0)
	for host := range existing {
		_, ok := tracker[host]
		if !ok {
			purge = append(purge, host)
		}
	}

	for _, each := range purge {
		delete(existing, each)
	}
}

func (proxy *Proxy) leastBusy() string {
	selectedHost := ``
	selectedHostLoad := 0

	hosts := proxy.hosts
	for host, load := range hosts {
		if selectedHost == `` {
			selectedHost = host
			selectedHostLoad = load
			continue
		}

		if selectedHostLoad > load {
			selectedHost = host
			selectedHostLoad = load
		}
	}
	return selectedHost
}

func (proxy *Proxy) HandleRequest(response Response, request *Request) {
	host := proxy.leastBusy()
	if host == `` {
		response.Error(
			http.StatusServiceUnavailable,
			errors.New(`DOWNSTREAM SERVICE UNAVAILABLE`),
		)
		return
	}

	downstream, err := url.Parse(host)
	if err != nil {
		response.Error(http.StatusServiceUnavailable, err)
		return
	}

	reverseProxy := httputil.ReverseProxy{
		Director: func(request *http.Request) {
			headers := request.Header
			request.Header = updateHeaders(
				request.Header,
				map[string][]string{
					XForwardedFor: {
						strings.Join(
							[]string{
								headers.Get(XForwardedFor),
								request.RemoteAddr,
							},
							`,`,
						),
					},
				},
			)
			request.URL = proxy.downstreamUrl(request.URL, downstream)
		},
	}

	native := *request
	response.compress = false
	reverseProxy.ServeHTTP(response, &native.Request)
}

func NewProxy(downstreamUrl func(url *url.URL, host *url.URL) *url.URL) *Proxy {
	return &Proxy{
		downstreamUrl: downstreamUrl,
		hosts:         make(map[string]int),
	}
}
