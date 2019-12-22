package av

import (
	"context"
	"net/http"
	"net/url"
)

// Connection is an interface that requests data from a server
type Connection interface {
	// Request creates an http Response from the given endpoint URL
	Request(ctx context.Context, endpoint *url.URL) (*http.Response, error)
}

type avConnection struct {
	copts connOptions
}

func defaultConnOptions() connOptions {
	return connOptions{
		client:  &http.Client{},
		host:    HostDefault,
		timeout: TimeoutDefault,
		rl:      NewRateLimiter(0, 0),
	}
}

// NewConnectionHost creates a new connection at the default Alpha Vantage host
func NewConnection(opts ...ConnOption) Connection {
	av := &avConnection{
		copts: defaultConnOptions(),
	}

	for _, opt := range opts {
		opt.apply(&av.copts)
	}

	return av
}

func (conn *avConnection) Client() *http.Client {
	return conn.copts.client
}

func (conn *avConnection) Host() string {
	return conn.copts.host
}

func (conn *avConnection) RateLimiter() *RateLimiter {
	return conn.copts.rl
}

// Request will make an HTTP GET request for the given endpoint from Alpha Vantage
func (conn *avConnection) Request(ctx context.Context, endpoint *url.URL) (*http.Response, error) {
	return conn.RateLimiter().Do(func() (*http.Response, error) {
		endpoint.Scheme = schemeHttps
		endpoint.Host = conn.Host()
		targetUrl := endpoint.String()

		req, err := http.NewRequest(http.MethodGet, targetUrl, nil)
		if err != nil {
			return nil, err
		}

		return conn.Client().Do(req.WithContext(ctx))
	})
}
