package tests

import (
	"context"
	"strconv"
	"testing"

	"github.com/szymongib/graphql-client/util"

	"github.com/stretchr/testify/require"

	"github.com/szymongib/graphql-client/test/schema"

	"github.com/stretchr/testify/assert"
	"github.com/szymongib/graphql-client/graphql"
)

func Test_E2e_Flow(t *testing.T) {
	// given
	defer resolver.ResetData()

	logger := func(s string) {
		assert.NotEmpty(t, s)
	}

	gqlClient := graphql.NewClient(apiAddress, graphql.WithLogger(logger))

	// create human
	tedHumanInput := humanInput("Ted", nil)

	var human schema.Human
	err := gqlClient.Mutate(context.Background(), "createHuman", graphql.OperationInput{"in": tedHumanInput}, &human)
	require.NoError(t, err)
	assert.NotEmpty(t, human.ID)
	assert.Equal(t, "Ted", human.Name)

	// query human
	var queriedHuman schema.Human
	err = gqlClient.Query(context.Background(), "human", graphql.OperationInput{"id": human.ID}, &queriedHuman)
	require.NoError(t, err)
	assert.Equal(t, human, queriedHuman)

	// add dogs to human
	dogsInput := []*schema.DogInput{
		dogInput("Dog1", util.IntPtr(1), []*schema.DistinguishingFeatureInput{{"Black spots", util.Float64Ptr(1)}}),
		dogInput("Dog2", util.IntPtr(2), []*schema.DistinguishingFeatureInput{{"Black spots", util.Float64Ptr(2)}}),
		dogInput("Dog3", util.IntPtr(3), []*schema.DistinguishingFeatureInput{{"Black spots", util.Float64Ptr(3)}}),
	}

	for i, dogInput := range dogsInput {
		var dog schema.Dog
		err = gqlClient.Mutate(context.Background(), "createDog", graphql.OperationInput{"humanID": human.ID, "in": dogInput}, &dog)
		require.NoError(t, err)
		assert.NotEmpty(t, dog.ID)
		assert.Equal(t, "Dog"+strconv.Itoa(i+1), dog.Name)
		assert.Equal(t, human.ID, dog.OwnerID)
		assert.Equal(t, util.IntPtr(i+1), dog.TailLength)
		assert.Equal(t, "Black spots", dog.DistinguishingFeatures[0].Description)
		assert.Equal(t, util.Float64Ptr(float64(i+1)), dog.DistinguishingFeatures[0].SpottingDifficulty)
	}

	// query dogs
	var dogs []*schema.Dog
	operation := graphql.Operation{
		Type:      graphql.Query,
		Name:      "dogs",
		Requested: &dogs,
		Input:     graphql.OperationInput{},
	}

	err = gqlClient.Run(context.Background(), operation, &dogs)
	require.NoError(t, err)
	assert.Equal(t, 3, len(dogs))
	for _, dog := range dogs {
		assert.Equal(t, human.ID, dog.OwnerID)
	}

	//// create another human
	//nedHumanInput := humanInput("Ned", nil)
	//req, err := graphql.NewMutationRequestFromData("createHuman", graphql.OperationInput{"in": nedHumanInput}, schema.Human{})
	//require.NoError(t, err)
	//
	//var nedHuman schema.Human
	//err = gqlClient.Execute(context.Background(), req, &nedHuman)
	//require.NoError(t, err)
	//assert.NotEmpty(t, nedHuman.ID)
	//assert.Equal(t, "Ned", nedHuman.Name)

}

func humanInput(name string, dogs []*schema.DogInput) *schema.HumanInput {
	return &schema.HumanInput{
		Name: name,
		Dogs: dogs,
	}
}

func dogInput(name string, tail *int, features []*schema.DistinguishingFeatureInput) *schema.DogInput {
	return &schema.DogInput{
		Name:                   name,
		TailLength:             tail,
		DistinguishingFeatures: features,
	}
}
