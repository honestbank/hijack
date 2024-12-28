package handlers

import (
	"fmt"

	"github.com/honestbank/hijack/v2"
	"github.com/honestbank/hijack/v2/request"
)

type handlers struct {
	items map[string]hijack.Hijacker
}

func (h *handlers) Set(operationName string, item func(r *request.GraphqlRequest) (string, error)) Manager {
	h.items[operationName] = item
	return h
}

func (h *handlers) Handle(r *request.GraphqlRequest) (string, error) {
	handler := h.items[r.OperationName]
	if handler == nil {
		return "", panic("no handler found for %s", r.OperationName)
	}

	return handler(r)
}

type Manager interface {
	Set(operationName string, item func(r *request.GraphqlRequest) (string, error)) Manager
	Handle(r *request.GraphqlRequest) (string, error)
}

func New() Manager {
	return &handlers{items: map[string]hijack.Hijacker{}}
}
