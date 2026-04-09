package service

import (
	"errors"
	"pcbook/pb"
	"pcbook/util"
	"sync"
)

// ErrAlreadyExists is returned when the store already contains a laptop with the same ID.
var ErrAlreadyExists = errors.New("laptop already exists")

// ErrNotFound is returned when the store does not contain a laptop with the given ID.
var ErrNotFound = errors.New("laptop not found")

// LaptopStore is an interface that defines the behavior for managing laptop data storage.
// It allows different implementations such as in-memory, database, or cloud storage.
type LaptopStore interface {
	// Save stores a new laptop in the system. Returns an error if the operation fails.
	Save(laptop *pb.Laptop) error
	// Get returns the laptop with the given ID.
	Get(id string) (*pb.Laptop, error)
}

// InMemoryLaptopStore is a thread-safe implementation of LaptopStore.
// It uses an internal map to store laptop information in the application's memory.
type InMemoryLaptopStore struct {
	mutex   sync.RWMutex
	laptops map[string]*pb.Laptop
}

// NewInMemoryLaptopStore creates and initializes a new InMemoryLaptopStore with an empty map.
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		laptops: make(map[string]*pb.Laptop),
	}
}

// Save attempts to persist a laptop message to the in-memory map.
// It returns ErrAlreadyExists if a laptop with the same unique identifier is already present.
// This operation is protected by a mutex lock to ensure safe concurrent access.
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.laptops[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other := util.DeepCopyLaptop(laptop)
	store.laptops[other.Id] = other
	return nil
}

// Get returns the laptop with the given ID.
func (store *InMemoryLaptopStore) Get(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.laptops[id]
	if laptop == nil {
		return nil, ErrNotFound
	}

	return util.DeepCopyLaptop(laptop), nil
}
