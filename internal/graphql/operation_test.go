package graphql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/hijack/internal/graphql"
)

const multipleGraphqlOperationDefinition = `
query {
  me {id}
}
mutation {
  doSomething {phone}
}
`
const graphqlOperationWithFragment = `
query GetUser($userId: ID!) {
  user(id: $userId) {
    id,
    name,
    isViewerFriend,
    profilePicture(size: 50)  {
      ...PictureFragment
    }
  }
}

fragment PictureFragment on Picture {
  uri,
  width,
  height
}
`
const invalidOperation = `
query GetUserId: ID!) {
  user(id: userId) {
    id,
    name,
    isViewerFriend,
  }
}
`

func TestGetOperationName(t *testing.T) {
	t.Run("one operation definition", func(t *testing.T) {
		t.Run("if there's more than 1 OperationDefinition, it returns error", func(t *testing.T) {
			_, err := graphql.GetOperationName(multipleGraphqlOperationDefinition)
			assert.EqualError(t, err, "expected exactly 1 OperationDefinition, found 2")
		})
		t.Run("allows fragments and only checks OperationDefinition", func(t *testing.T) {
			name, err := graphql.GetOperationName(graphqlOperationWithFragment)
			assert.NoError(t, err)
			assert.Equal(t, "GetUser", name)
		})
	})
	t.Run("returns error if syntax is wrong", func(t *testing.T) {
		_, err := graphql.GetOperationName(invalidOperation)
		assert.Error(t, err)
	})
	t.Run("returns error if the operation name is empty", func(t *testing.T) {
		_, err := graphql.GetOperationName(`query{me{id}}`)
		assert.EqualError(t, err, "an operation must have a name")
	})
}
