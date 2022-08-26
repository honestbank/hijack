package handlers

import (
	"fmt"

	"github.com/honestbank/hijack"
)

type handlers struct {
	items map[string]hijack.Hijacker
}

func (h *handlers) Set(operationName string, item func(operation string) (string, error)) Manager {
	h.items[operationName] = item
	return h
}

func (h *handlers) Handle(operationName string) (string, error) {
	handler := h.items[operationName]
	if handler == nil {
		return "", fmt.Errorf("no handler found for %s", operationName)
	}

	return handler(operationName)
}

type Manager interface {
	Set(operationName string, item func(operationName string) (string, error)) Manager
	Handle(operationName string) (string, error)
}

func New() Manager {
	return &handlers{items: map[string]hijack.Hijacker{}}
}
