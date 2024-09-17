package transport

import (
	"net/http"
)

type transport struct {
	h                 Handler
	host              string
	originalTransport http.RoundTripper
}

type Handler func(r *http.Request) (*http.Response, error)

func (t *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.host != request.Host {
		return t.originalTransport.RoundTrip(request)
	}
	return t.h(request)
}

func Hijack(handler Handler, host string) func() {
	currentTransport := http.DefaultTransport
	http.DefaultTransport = &transport{h: handler, host: host, originalTransport: currentTransport}

	return func() {
		http.DefaultTransport = currentTransport
	}
}
