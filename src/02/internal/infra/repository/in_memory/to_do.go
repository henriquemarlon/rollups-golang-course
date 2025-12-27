package in_memory

import (
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/domain"
)

func (r *InMemoryRepository) CreateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	input.Id = r.NextID
	r.NextID++
	r.Db[input.Id] = input
	return input, nil
}

func (r *InMemoryRepository) FindAllToDos() ([]*domain.ToDo, error) {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	var todos []*domain.ToDo
	for _, todo := range r.Db {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *InMemoryRepository) UpdateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	todo, exists := r.Db[input.Id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	todo.Title = input.Title
	todo.Description = input.Description
	todo.Completed = input.Completed

	r.Db[input.Id] = todo

	return todo, nil
}

func (r *InMemoryRepository) DeleteToDo(id uint) error {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	_, exists := r.Db[id]
	if !exists {
		return domain.ErrNotFound
	}

	delete(r.Db, id)
	return nil
}
