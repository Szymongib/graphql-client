package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/szymongib/graphql-client/graphql"
)

func Test_Errors(t *testing.T) {

	t.Run("should return proper graphql error", func(t *testing.T) {
		client := graphql.NewClient(apiAddress)

		var response string
		err := client.Execute(context.Background(), graphql.NewRequestRaw("query{ errorsQuery }"), &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "error you requested")

		err = client.Query(context.Background(), "errorsQuery", nil, &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "error you requested")

		err = client.Mutate(context.Background(), "errorsMutation", nil, &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "error you requested")
	})

	t.Run("should return error when server responded with non 200 code without GQL errors", func(t *testing.T) {
		client := graphql.NewClient(errorsAddress + "/noGQL")

		var response interface{}
		err := client.Execute(context.Background(), graphql.NewRequestRaw(""), &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "Response body:")
		assert.Contains(t, err.Error(), "unexpected response status")
	})

	t.Run("should return error when server responded with 200 code with non GQL response", func(t *testing.T) {
		client := graphql.NewClient(nonGQLAddress)

		var response interface{}
		err := client.Execute(context.Background(), graphql.NewRequestRaw(""), &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "failed to decode response body")
	})

	t.Run("should return error on invalid operation", func(t *testing.T) {
		client := graphql.NewClient(apiAddress)

		var response interface{}
		err := client.Execute(context.Background(), graphql.NewRequestRaw("query { result: invalidOperation }"), &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "Cannot query field")
	})

	t.Run("should return error on invalid query syntax", func(t *testing.T) {
		client := graphql.NewClient(apiAddress)

		var response interface{}
		err := client.Execute(context.Background(), graphql.NewRequestRaw("query { result: dogs }}"), &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "Unexpected }")
	})

	t.Run("should return error if context canceled", func(t *testing.T) {
		client := graphql.NewClient(apiAddress)

		ctx, _ := context.WithTimeout(context.Background(), 1*time.Millisecond)
		time.Sleep(1 * time.Second)

		var response interface{}
		err := client.Execute(ctx, graphql.NewRequestRaw("query { result: dogs }"), &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "context deadline exceeded")
	})

	t.Run("should return error if invalid input to map", func(t *testing.T) {
		client := graphql.NewClient(apiAddress)

		var response interface{}
		err := client.Query(context.Background(), "dogs", graphql.OperationInput{"in": struct {
			MyMap map[interface{}]interface{}
		}{}}, &response)
		require.Error(t, err)
		t.Log(err.Error())
		assert.Contains(t, err.Error(), "failed to parse to GQL input")
	})
}
