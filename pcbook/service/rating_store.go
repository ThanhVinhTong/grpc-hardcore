package service

import (
	"context"
	"sync"
)

type RatingStore interface {
	Add(ctx context.Context, laptopID string, score float64) (*Rating, error)
}

type Rating struct {
	RatedCount uint32
	Sum        float64
}

type InMemoryRatingStore struct {
	mutex  sync.RWMutex
	rating map[string]*Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

func (store *InMemoryRatingStore) Add(ctx context.Context, laptopID string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopID]
	if rating == nil {
		rating = &Rating{
			RatedCount: 1,
			Sum:        score,
		}
	} else {
		rating.RatedCount++
		rating.Sum += score
	}

	store.rating[laptopID] = rating
	return rating, nil
}
