package tests

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/szymongib/graphql-client/test/schema"
	"github.com/szymongib/graphql-client/util"

	"github.com/szymongib/graphql-client/graphql"
)

func Test_Headers(t *testing.T) {
	defer resolver.ResetData()

	gqlClient := graphql.NewClient(apiAddress, graphql.WithLogger(func(s string) {
		fmt.Println(s)
	}))

	headers := http.Header{
		"Test":           {"val1", "val2"},
		"Another-Header": {"header-value"},
	}

	expectedHeaders := []*schema.Header{
		{
			Name:   "Test",
			Values: []*string{util.StringPtr("val1"), util.StringPtr("val2")},
		},
		{
			Name:   "Another-Header",
			Values: []*string{util.StringPtr("header-value")},
		},
	}

	t.Run("should return headers with raw query", func(t *testing.T) {
		request := graphql.NewRequestRaw("query { result: headersQuery() { name values } }", headers)

		var actualHeaders struct {
			Result []*schema.Header `json:"result"`
		}
		err := gqlClient.Execute(context.Background(), request, &actualHeaders)
		require.NoError(t, err)
		assert.True(t, containsHeaders(actualHeaders.Result, expectedHeaders))
	})

	t.Run("should return headers with raw mutation", func(t *testing.T) {
		request := graphql.NewRequestRaw("mutation { result: headersMutation() { name values } }", headers)

		var actualHeaders struct {
			Result []*schema.Header `json:"result"`
		}
		err := gqlClient.Execute(context.Background(), request, &actualHeaders)
		require.NoError(t, err)
		assert.True(t, containsHeaders(actualHeaders.Result, expectedHeaders))
	})

	t.Run("should return headers with mapped query", func(t *testing.T) {
		var actualHeaders []*schema.Header
		err := gqlClient.Query(context.Background(), "headersQuery", nil, &actualHeaders, headers)
		require.NoError(t, err)
		assert.True(t, containsHeaders(actualHeaders, expectedHeaders))
	})

	t.Run("should return headers with mapped mutation", func(t *testing.T) {
		var actualHeaders []*schema.Header
		err := gqlClient.Mutate(context.Background(), "headersMutation", nil, &actualHeaders, headers)
		require.NoError(t, err)
		assert.True(t, containsHeaders(actualHeaders, expectedHeaders))
	})

	t.Run("should return headers with mapped query operation", func(t *testing.T) {
		var actualHeaders []*schema.Header

		operation := graphql.Operation{
			Type: graphql.Query,
			Name: "headersQuery",
		}

		err := gqlClient.Run(context.Background(), operation, &actualHeaders, headers)
		require.NoError(t, err)
		assert.True(t, containsHeaders(actualHeaders, expectedHeaders))
	})

	t.Run("should return headers with mapped mutation operation", func(t *testing.T) {
		var actualHeaders []*schema.Header

		operation := graphql.Operation{
			Type: graphql.Mutation,
			Name: "headersMutation",
		}

		err := gqlClient.Run(context.Background(), operation, &actualHeaders, headers)
		require.NoError(t, err)
		assert.True(t, containsHeaders(actualHeaders, expectedHeaders))
	})

}

func containsHeaders(all, contains []*schema.Header) bool {
	for _, header := range contains {
		if !containsHeader(header, all) {
			return false
		}
	}

	return true
}

func containsHeader(element *schema.Header, collection []*schema.Header) bool {
	for _, e := range collection {
		if reflect.DeepEqual(element, e) {
			return true
		}
	}
	return false
}
