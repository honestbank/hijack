package request_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/internal/request"
)

const graphqlOperationWithFragment = `query GetUser($userId: ID!) { user(id: $userId) { id  }}`

func TestGetOperation(t *testing.T) {
	t.Run("for get request with empty body, it returns EOF error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(""))
		assert.NoError(t, err)
		_, err = request.GetOperation(req)
		assert.EqualError(t, err, "EOF")
	})
	t.Run("if it doesn't find operation name in body, it returns an error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"query": ""}`))
		assert.NoError(t, err)
		_, err = request.GetOperation(req)
		assert.EqualError(t, err, "expected exactly 1 OperationDefinition, found 0")
	})
	t.Run("returns valid operation name from request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{"query": "%s"}`, graphqlOperationWithFragment)))
		assert.NoError(t, err)
		name, err := request.GetOperation(req)
		assert.NoError(t, err)
		assert.Equal(t, "GetUser", name)
	})
}
