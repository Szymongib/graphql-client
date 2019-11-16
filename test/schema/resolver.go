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

type Resolver struct {
	HumansDb []*Human
	DogsDb   []*Dog
}

func (r *Resolver) Mutation() MutationResolver {
	return r
}

func (r *Resolver) Query() QueryResolver {
	return r
}

func (r *Resolver) ResetData() {
	r.DogsDb = []*Dog{}
	r.HumansDb = []*Human{}
}

// Queries

func (r *Resolver) Humans(ctx context.Context) ([]*Human, error) {
	return r.HumansDb, nil
}

func (r *Resolver) Dogs(ctx context.Context) ([]*Dog, error) {
	return r.DogsDb, nil
}

func (r *Resolver) Human(ctx context.Context, id string) (*Human, error) {
	return r.getHumanById(id)
}

func (r *Resolver) Dog(ctx context.Context, id string) (*Dog, error) {
	return r.getDogById(id)
}

func (r *Resolver) HeadersQuery(ctx context.Context) ([]*Header, error) {
	return getHeaders(ctx), nil
}

func (r *Resolver) ErrorsQuery(ctx context.Context) (string, error) {
	return "", fmt.Errorf("error you requested")
}

// Mutations

func (r *Resolver) ErrorsMutation(ctx context.Context) (string, error) {
	return "", fmt.Errorf("error you requested")
}

func (r *Resolver) HeadersMutation(ctx context.Context) ([]*Header, error) {
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

func (r *Resolver) CreateHuman(ctx context.Context, in HumanInput) (*Human, error) {
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

func (r *Resolver) CreateDog(ctx context.Context, humanID string, in DogInput) (*Dog, error) {
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
			Description:        feature.Description,
			SpottingDifficulty: feature.SpottingDifficulty,
		}

		newFeatures = append(newFeatures, &f)
	}

	return newFeatures
}

func (r *Resolver) getDogById(dogId string) (*Dog, error) {
	for _, d := range r.DogsDb {
		if d.ID == dogId {
			return d, nil
		}
	}

	return nil, fmt.Errorf("dog with ID %s not found", dogId)
}

func (r *Resolver) getHumanById(humanId string) (*Human, error) {
	for _, h := range r.HumansDb {
		if h.ID == humanId {
			return h, nil
		}
	}

	return nil, fmt.Errorf("human with ID %s not found", humanId)
}
