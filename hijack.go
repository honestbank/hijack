package hijack

import (
	"io"
	"net/http"
	"strings"

	"github.com/honestbank/hijack/internal/request"
	"github.com/honestbank/hijack/internal/transport"
)

type Hijacker func(operationName string) (string, error)

// Start a hijacking session where all requests going to `host` will be mocked.
func Start(hijacker Hijacker, host string) func() {
	return transport.Hijack(func(r *http.Request) (*http.Response, error) {
		operationName, err := request.GetOperation(r)
		if err != nil {
			return nil, err
		}
		response, err := hijacker(operationName)
		if err != nil {
			return nil, err
		}
		return &http.Response{
			Body: io.NopCloser(strings.NewReader(response)),
		}, nil
	}, host)
}
