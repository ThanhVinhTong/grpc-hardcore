package service

import (
	"context"
	"errors"
	"log"
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
	// Search returns a list of laptops that match the given filter.
	Search(ctx context.Context, filter *pb.Filter, found func(*pb.Laptop) error) error
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

// Search searches for laptops with filter, returns one by one via the found function
func (store *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *pb.Filter,
	found func(laptop *pb.Laptop) error,
) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.laptops {
		// heavy duties simulation
		// time.Sleep(time.Second)
		// log.Printf("laptop with id %s said ehehehe", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Println("context cancelled or deadline exceeded")
			return errors.New("context is cancelled")
		}

		if isQualified(filter, laptop) {
			laptop_copy := util.DeepCopyLaptop(laptop)

			err := found(laptop_copy)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3
	case pb.Memory_KILOBYTE:
		return value << 13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}
