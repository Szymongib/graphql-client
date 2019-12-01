package schema

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	HeadersContextKey = "Headers"
)

func NewResolver() Resolver {
	db := &Database{
		HumansDb: []*Human{},
		DogsDb:   []*Dog{},
	}

	return Resolver{
		Database:         db,
		humanResolver:    humanResolver{Database: db},
		queryResolver:    queryResolver{Database: db},
		mutationResolver: mutationResolver{Database: db},
	}
}

type Database struct {
	HumansDb []*Human
	DogsDb   []*Dog
}

type Resolver struct {
	*Database
	humanResolver    humanResolver
	queryResolver    queryResolver
	mutationResolver mutationResolver
}

type queryResolver struct {
	*Database
}

type mutationResolver struct {
	*Database
}

type humanResolver struct {
	*Database
}

func (r *Resolver) Mutation() MutationResolver {
	return &r.mutationResolver
}

func (r *Resolver) Query() QueryResolver {
	return &r.queryResolver
}

func (r *Resolver) Human() HumanResolver {
	return &r.humanResolver
}

func (r *Resolver) ResetData() {
	r.DogsDb = []*Dog{}
	r.HumansDb = []*Human{}
}

// Queries

func (r *queryResolver) Humans(ctx context.Context) ([]*Human, error) {
	return r.HumansDb, nil
}

func (r *humanResolver) Dogs(ctx context.Context, obj *Human, limit *int, offset *int, filters []*FeatureFilterInput) ([]*Dog, error) {
	if len(r.DogsDb) == 0 {
		return r.DogsDb, nil
	}

	var dogs []*Dog

	if filters != nil {
		dogs = filterDogs(r.DogsDb, obj.ID, filters)
	} else {
		dogs = r.DogsDb
	}

	dogsOffset := 0
	if offset != nil {
		dogsOffset = *offset
		if dogsOffset >= len(dogs) {
			dogsOffset = len(dogs) - 1
		}
	}

	end := len(dogs)
	if limit != nil {
		end := dogsOffset + *limit
		if end >= len(dogs) {
			end = len(dogs) - 1
		}
	}

	return dogs[dogsOffset:end], nil
}

func filterDogs(dogs []*Dog, humanId string, filters []*FeatureFilterInput) []*Dog {
	var filteredDogs []*Dog

	for _, f := range filters {
		for _, d := range dogs {
			if d.OwnerID == humanId && FeaturesSatisfiesFilter(d.DistinguishingFeatures, *f) {
				filteredDogs = append(filteredDogs, d)
			}
		}
	}

	return filteredDogs
}

func (r *queryResolver) Human(ctx context.Context, id string) (*Human, error) {
	return r.getHumanById(id)
}

func (r *queryResolver) Dog(ctx context.Context, id string) (*Dog, error) {
	return r.getDogById(id)
}

func (r *queryResolver) Dogs(ctx context.Context) ([]*Dog, error) {
	return r.DogsDb, nil
}

func (r *queryResolver) HeadersQuery(ctx context.Context) ([]*Header, error) {
	return getHeaders(ctx), nil
}

func (r *queryResolver) ErrorsQuery(ctx context.Context) (string, error) {
	return "", fmt.Errorf("error you requested")
}

// Mutations

func (r *mutationResolver) ErrorsMutation(ctx context.Context) (string, error) {
	return "", fmt.Errorf("error you requested")
}

func (r *mutationResolver) HeadersMutation(ctx context.Context) ([]*Header, error) {
	return getHeaders(ctx), nil
}

func getHeaders(ctx context.Context) []*Header {
	headersCtx := ctx.Value(HeadersContextKey)

	httpHeaders, ok := headersCtx.(http.Header)
	if !ok {
		return []*Header{}
	}

	headers := make([]*Header, 0, len(httpHeaders))

	for k, v := range httpHeaders {
		headers = append(headers, &Header{
			Name:   k,
			Values: strPtrSlice(v),
		})
	}

	return headers
}

func strPtrSlice(slice []string) []*string {
	strPtrs := make([]*string, len(slice))

	for i, _ := range slice {
		strPtrs[i] = &slice[i]
	}

	return strPtrs
}

func (r *mutationResolver) CreateHuman(ctx context.Context, in HumanInput) (*Human, error) {
	humanId := uuid.New().String()

	dogs := make([]*Dog, 0, len(in.Dogs))

	for _, dog := range in.Dogs {
		dogId := uuid.New().String()
		newDog := &Dog{
			ID:                     dogId,
			Name:                   dog.Name,
			OwnerID:                humanId,
			TailLength:             dog.TailLength,
			DistinguishingFeatures: newDistinguishingFeatures(dog.DistinguishingFeatures),
		}
		dogs = append(dogs, newDog)
	}

	human := &Human{
		ID:   humanId,
		Name: in.Name,
		Dogs: dogs,
	}

	r.HumansDb = append(r.HumansDb, human)

	return human, nil
}

func (r *mutationResolver) CreateDog(ctx context.Context, humanID string, in DogInput) (*Dog, error) {
	human, err := r.getHumanById(humanID)
	if err != nil {
		return nil, err
	}

	newDogId := uuid.New().String()

	newDog := &Dog{
		ID:                     newDogId,
		Name:                   in.Name,
		OwnerID:                human.ID,
		TailLength:             in.TailLength,
		DistinguishingFeatures: newDistinguishingFeatures(in.DistinguishingFeatures),
	}

	r.DogsDb = append(r.DogsDb, newDog)

	return newDog, nil
}

func newDistinguishingFeatures(input []*DistinguishingFeatureInput) []*DistinguishingFeature {
	newFeatures := make([]*DistinguishingFeature, 0, len(input))

	for _, feature := range input {
		f := DistinguishingFeature{
			Name:        feature.Name,
			Description: feature.Description,
			Intensity:   feature.Intensity,
		}

		newFeatures = append(newFeatures, &f)
	}

	return newFeatures
}

func (r *Database) getDogById(dogId string) (*Dog, error) {
	for _, d := range r.DogsDb {
		if d.ID == dogId {
			return d, nil
		}
	}

	return nil, fmt.Errorf("dog with ID %s not found", dogId)
}

func (r *Database) getHumanById(humanId string) (*Human, error) {
	for _, h := range r.HumansDb {
		if h.ID == humanId {
			return h, nil
		}
	}

	return nil, fmt.Errorf("human with ID %s not found", humanId)
}
