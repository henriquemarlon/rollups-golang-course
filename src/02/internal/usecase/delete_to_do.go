package usecase

import "github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository"

type DeleteToDoInputDTO struct {
	Id uint `json:"id" validate:"required"`
}

type DeleteToDoUseCase struct {
	ToDoRepository repository.ToDoRepository
}

func NewDeleteToDoUseCase(todoRepository repository.ToDoRepository) *DeleteToDoUseCase {
	return &DeleteToDoUseCase{
		ToDoRepository: todoRepository,
	}
}

func (u *DeleteToDoUseCase) Execute(input *DeleteToDoInputDTO) error {
	return u.ToDoRepository.DeleteToDo(input.Id)
}
