package hijack_test

import (
	"context"
	"errors"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/v2"
	"github.com/honestbank/hijack/v2/internal/request"
)

type Response struct {
	Me string `json:"me"`
}

func TestHijack(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		restoreFunction := hijack.Start(func(r *request.GraphqlRequest) (string, error) {
			if r.OperationName == "GetUser" {
				return `{"data": {"me": "something"}}`, nil
			}
			return "", nil
		}, "localhost:8000")
		r := graphql.NewRequest("query GetUser {me}")
		client := graphql.NewClient("http://localhost:8000/graphql")
		resp := Response{}
		err := client.Run(context.Background(), r, &resp)
		assert.NoError(t, err)
		assert.Equal(t, "something", resp.Me)
		// make your graphql calls
		restoreFunction()
		// now it'll fire real requests
	})
	t.Run("fire invalid query, receive an error", func(t *testing.T) {
		defer hijack.Start(func(r *request.GraphqlRequest) (string, error) {
			if r.OperationName == "GetUser" {
				return `{"data": {"me": "something"}}`, nil
			}
			return "", nil
		}, "localhost:8000")()
		r := graphql.NewRequest("query GetUser me}")
		client := graphql.NewClient("http://localhost:8000/graphql")
		resp := Response{}
		err := client.Run(context.Background(), r, &resp)
		assert.Error(t, err)
	})
	t.Run("hijacker can return errors", func(t *testing.T) {
		defer hijack.Start(func(r *request.GraphqlRequest) (string, error) {
			if r.OperationName == "GetUser" {
				return "", errors.New("some error")
			}
			return "", nil
		}, "localhost:8000")()
		r := graphql.NewRequest("query GetUser {me}")
		client := graphql.NewClient("http://localhost:8000/graphql")
		resp := Response{}
		err := client.Run(context.Background(), r, &resp)
		assert.Error(t, err)
	})
}
