package transport_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/v2/internal/transport"
)

func TestHijack(t *testing.T) {
	t.Run("can hijack default transport and restore", func(t *testing.T) {
		restoreFn := transport.Hijack(func(request *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("error for %s", request.URL)
		}, "localhost")
		res, err := http.Get("http://localhost/something")
		assert.Nil(t, res)
		assert.EqualError(t, err, "Get \"http://localhost/something\": error for http://localhost/something")
		restoreFn()
	})
	t.Run("can return valid response", func(t *testing.T) {
		defer transport.Hijack(func(request *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 444}, nil
		}, "localhost")()
		res, err := http.Get("http://localhost/something")
		assert.NoError(t, err)
		assert.Equal(t, 444, res.StatusCode)
	})
	t.Run("still forwards request when host is not matched", func(t *testing.T) {
		called := false
		srv := http.Server{Addr: ":3000"}
		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		}))
		go func() {
			_ = srv.ListenAndServe()
		}()
		defer transport.Hijack(func(request *http.Request) (*http.Response, error) {
			called = true

			return &http.Response{StatusCode: 444}, nil
		}, "localhost")()
		_, err := http.Get("http://localhost:3000")
		_ = srv.Shutdown(context.Background())
		assert.NoError(t, err)
		assert.False(t, called)
	})
}
