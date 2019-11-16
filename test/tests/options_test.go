package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/szymongib/graphql-client/test/schema"

	"github.com/szymongib/graphql-client/graphql"
)

func Test_Options(t *testing.T) {
	noTimeoutHttpClient := &http.Client{
		Timeout: 1 * time.Nanosecond,
	}

	client := graphql.NewClient(apiAddress, graphql.WithHTTPClient(noTimeoutHttpClient))

	var dogs []*schema.Dog
	err := client.Query(context.Background(), "dogs", nil, &dogs)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "Client.Timeout")
}
