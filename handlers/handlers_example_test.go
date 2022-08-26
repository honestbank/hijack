package handlers_test

import (
	"fmt"

	"github.com/honestbank/hijack/handlers"
)

func ExampleNew() {
	cb := func(operationName string) (string, error) {
		return "hello from " + operationName, nil
	}
	h := handlers.New().
		Set("Operation1", cb).
		Set("Operation2", cb).
		Set("Operation3", cb)
	res, _ := h.Handle("Operation2")
	fmt.Println(res)
	// Output: hello from Operation2
}
