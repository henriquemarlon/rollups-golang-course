package in_memory

import (
	"sync"

	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/domain"
)

type InMemoryRepository struct {
	Db     map[uint]*domain.ToDo
	Mutex  *sync.RWMutex
	NextID uint
}

func (r *InMemoryRepository) Close() error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Db = make(map[uint]*domain.ToDo)
	r.NextID = 1
	return nil
}

func NewInMemoryRepository() (*InMemoryRepository, error) {
	return &InMemoryRepository{
		Db:     make(map[uint]*domain.ToDo),
		Mutex:  &sync.RWMutex{},
		NextID: 1,
	}, nil
}
