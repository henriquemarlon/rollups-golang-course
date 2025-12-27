package usecase

import (
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/pkg/rollups"
)

type CreateToDoInputDTO struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type CreateToDoOutputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   uint64 `json:"created_at"`
}

type CreateToDoUseCase struct {
	ToDoRepository repository.ToDoRepository
}

func NewCreateToDoUseCase(todoRepository repository.ToDoRepository) *CreateToDoUseCase {
	return &CreateToDoUseCase{
		ToDoRepository: todoRepository,
	}
}

func (u *CreateToDoUseCase) Execute(input *CreateToDoInputDTO, metadata rollups.Metadata) (*CreateToDoOutputDTO, error) {
	res, err := domain.NewToDo(input.Title, input.Description, metadata.BlockTimestamp)
	if err != nil {
		return nil, err
	}

	res, err = u.ToDoRepository.CreateToDo(res)
	if err != nil {
		return nil, err
	}

	return &CreateToDoOutputDTO{
		Id:          res.Id,
		Title:       res.Title,
		Description: res.Description,
		Completed:   res.Completed,
		CreatedAt:   res.CreatedAt,
	}, nil
}
