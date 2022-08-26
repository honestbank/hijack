package handlers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/handlers"
)

func mockHandler(operationName string) (string, error) {
	return operationName, nil
}

func TestNew(t *testing.T) {
	h := handlers.New()
	h.Set("1", mockHandler).
		Set("2", mockHandler).
		Set("3", mockHandler).
		Set("4", mockHandler)
	res, _ := h.Handle("1")
	assert.Equal(t, "1", res)
	res, _ = h.Handle("2")
	assert.Equal(t, "2", res)
	res, _ = h.Handle("3")
	assert.Equal(t, "3", res)
	res, _ = h.Handle("4")
	assert.Equal(t, "4", res)
	res, err := h.Handle("5")
	assert.Equal(t, "", res)
	assert.EqualError(t, err, "no handler found for 5")
}
