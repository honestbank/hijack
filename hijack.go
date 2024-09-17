package hijack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/honestbank/hijack/v2/internal/request"
	"github.com/honestbank/hijack/v2/internal/transport"
)

func GetVariableAs[T any](r *request.GraphqlRequest, name string) (T, error) {
	var output T
	variable, ok := r.Variables[name]
	if !ok {
		return output, fmt.Errorf("variable %s not found", name)
	}
	err := json.Unmarshal(variable, &output)
	if err != nil {
		return output, err
	}
	return output, nil
}

type Hijacker func(o *request.GraphqlRequest) (string, error)

// Start a hijacking session where all requests going to `host` will be mocked.
func Start(hijacker Hijacker, host string) func() {
	return transport.Hijack(func(r *http.Request) (*http.Response, error) {
		parsedRequest, err := request.Parse(r)
		if err != nil {
			return nil, err
		}
		response, err := hijacker(parsedRequest)
		if err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(response)),
		}, nil
	}, host)
}
