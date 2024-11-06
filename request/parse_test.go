package request_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/v2/request"
)

const graphqlOperationWithFragment = `query GetUser($userId: ID!) { user(id: $userId) { id  }}`
const graphqlOperationVariables = `{"userId": 123}`

func TestParse(t *testing.T) {
	t.Run("for get request with empty body, it returns EOF error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(""))
		assert.NoError(t, err)
		_, err = request.Parse(req)
		assert.EqualError(t, err, "EOF")
	})
	t.Run("if it doesn't find operation name in body, it returns an error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"query": ""}`))
		assert.NoError(t, err)
		_, err = request.Parse(req)
		assert.EqualError(t, err, "expected exactly 1 OperationDefinition, found 0")
	})
	t.Run("returns valid operation name from query field", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{"query": "%s"}`, graphqlOperationWithFragment)))
		assert.NoError(t, err)
		r, err := request.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, "GetUser", r.OperationName)
	})
	t.Run("returns valid variables", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{"query": "%s", "variables": %s}`, graphqlOperationWithFragment, graphqlOperationVariables)))
		assert.NoError(t, err)
		r, err := request.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, "GetUser", r.OperationName)
	})
	t.Run("hydrates valid request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(fmt.Sprintf(`{"query": "%s", "variables": %s}`, graphqlOperationWithFragment, graphqlOperationVariables)))
		req.Header.Set("x-user-id", "123")
		assert.NoError(t, err)
		r, err := request.Parse(req)
		assert.NoError(t, err)
		assert.Equal(t, "123", r.OriginalRequest.Header.Get("x-user-id"))
	})
}
