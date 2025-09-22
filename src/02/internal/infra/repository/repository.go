package repository

import "github.com/henriquemarlon/cartesi-golang-series/src/02/internal/domain"

type ToDoRepository interface {
	CreateToDo(toDo *domain.ToDo) (*domain.ToDo, error)
	FindAllToDos() ([]*domain.ToDo, error)
	UpdateToDo(toDo *domain.ToDo) (*domain.ToDo, error)
	DeleteToDo(id uint) error
}

type Repository interface {
	ToDoRepository
	Close() error
}
