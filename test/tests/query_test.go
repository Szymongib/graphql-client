package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/szymongib/graphql-client/test/schema"

	"github.com/szymongib/graphql-client/graphql"
)

// TODO: nested queries test

func Test_Query(t *testing.T) {
	gqlClient := graphql.NewClient(apiAddress)

	dogID := uuid.New().String()

	resolver.DogsDb = []*schema.Dog{
		{
			ID:                     dogID,
			Name:                   "test",
			OwnerID:                uuid.New().String(),
			TailLength:             nil,
			DistinguishingFeatures: nil,
		},
	}

	t.Run("query dogs", func(t *testing.T) {
		var dogs []*schema.Dog
		err := gqlClient.Query(context.Background(), "dogs", nil, &dogs)
		require.NoError(t, err)

		assert.Equal(t, 1, len(dogs))
	})

	t.Run("query dog with ID", func(t *testing.T) {
		input := graphql.OperationInput{
			"id": dogID,
		}

		var dog schema.Dog
		err := gqlClient.Query(context.Background(), "dog", input, &dog)
		require.NoError(t, err)

		assert.Equal(t, dogID, dog.ID)
		assert.Equal(t, "test", dog.Name)
	})

}
