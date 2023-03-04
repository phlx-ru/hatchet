package transport

import (
	"context"
	"net"
	"net/http"
	"time"
)

const (
	defaultTimout                = 30 * time.Second
	defaultKeepAlive             = 30 * time.Second
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTLSHandShakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
	defaultMaxIdleConnsPerHost   = 255
)

type HTTPClient interface {
	Do(r *http.Request) (*http.Response, error)
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimout,
		Transport: &http.Transport{
			DialContext:           dialContext(defaultTimout, defaultKeepAlive),
			MaxIdleConnsPerHost:   defaultMaxIdleConnsPerHost,
			IdleConnTimeout:       defaultIdleConnTimeout,
			TLSHandshakeTimeout:   defaultTLSHandShakeTimeout,
			ExpectContinueTimeout: defaultExpectContinueTimeout,
		},
	}
}

func dialContext(timeout, keepAlive time.Duration) func(ctx context.Context, network, address string) (
	net.Conn,
	error,
) {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: keepAlive,
	}
	return dialer.DialContext
}
