package usecase

import (
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/pkg/rollups"
)

type UpdateToDoInputDTO struct {
	Id          uint   `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Completed   bool   `json:"completed" validate:"required"`
}

type UpdateToDoOutputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   uint64 `json:"created_at"`
	UpdatedAt   uint64 `json:"updated_at"`
}

type UpdateToDoUseCase struct {
	ToDoRepository repository.ToDoRepository
}

func NewUpdateToDoUseCase(todoRepository repository.ToDoRepository) *UpdateToDoUseCase {
	return &UpdateToDoUseCase{
		ToDoRepository: todoRepository,
	}
}

func (u *UpdateToDoUseCase) Execute(input *UpdateToDoInputDTO, metadata rollups.Metadata) (*UpdateToDoOutputDTO, error) {
	res, err := u.ToDoRepository.UpdateToDo(&domain.ToDo{
		Id:          input.Id,
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
		UpdatedAt:   metadata.BlockTimestamp,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateToDoOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Completed:   res.Completed,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}, nil
}
