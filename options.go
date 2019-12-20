package av

import (
	"net/http"
	"time"
)

type connOptions struct {
	client  *http.Client
	host    string
	timeout time.Duration
	rl      *RateLimiter
}

type ConnOption interface {
	apply(*connOptions)
}

// funcConnOption wraps a function that modifies connOptions into an
// implementation of the ConnOption interface.
type funcConnOption struct {
	f func(*connOptions)
}

func (fdo *funcConnOption) apply(do *connOptions) {
	fdo.f(do)
}

func newFuncConnOption(f func(*connOptions)) *funcConnOption {
	return &funcConnOption{
		f: f,
	}
}

func WithHost(host string) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.host = host
	})
}

func WithRateLimiter(rl *RateLimiter) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.rl = rl
	})
}

func WithTimeout(timeout time.Duration) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		if o.client == nil {
			o.client = &http.Client{}
		}
		o.client.Timeout = timeout
	})
}

func WithHTTPClient(client *http.Client) ConnOption {
	return newFuncConnOption(func(o *connOptions) {
		o.client = client
	})
}

type ClientOption interface {
	apply(*clientOptions)
}

type clientOptions struct {
	apiKey string
	conn   Connection
}

// funcClientOption wraps a function that modifies connOptions into an
// implementation of the ClientOption interface.
type funcClientOption struct {
	f func(*clientOptions)
}

func (fdo *funcClientOption) apply(do *clientOptions) {
	fdo.f(do)
}

func newFuncClientOption(f func(*clientOptions)) *funcClientOption {
	return &funcClientOption{
		f: f,
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return newFuncClientOption(func(o *clientOptions) {
		o.apiKey = apiKey
	})
}

func WithConnection(conn Connection) ClientOption {
	return newFuncClientOption(func(o *clientOptions) {
		o.conn = conn
	})
}
