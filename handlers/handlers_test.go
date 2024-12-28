package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/v2/handlers"
	"github.com/honestbank/hijack/v2/request"
)

func mockHandler(r *request.GraphqlRequest) (string, error) {
	return r.OperationName, nil
}

func TestNew(t *testing.T) {
	h := handlers.New()
	h.Set("1", mockHandler).
		Set("2", mockHandler).
		Set("3", mockHandler).
		Set("4", mockHandler)
	req := request.GraphqlRequest{
		OperationName: "1",
	}
	res, _ := h.Handle(&req)
	assert.Equal(t, "1", res)

	req.OperationName = "2"
	res, _ = h.Handle(&req)
	assert.Equal(t, "2", res)

	req.OperationName = "3"
	res, _ = h.Handle(&req)
	assert.Equal(t, "3", res)

	req.OperationName = "4"
	res, _ = h.Handle(&req)
	assert.Equal(t, "4", res)

	req.OperationName = "5"
	assert.Panics(t, func() {
		_, _ = h.Handle(&req)
	})
}
