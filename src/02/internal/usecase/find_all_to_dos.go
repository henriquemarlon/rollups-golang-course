package usecase

import "github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository"

type FindToDoOutputDTO struct {
	Id          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   uint64 `json:"created_at"`
	UpdatedAt   uint64 `json:"updated_at"`
}

type FindAllToDosOutputDTO []*FindToDoOutputDTO

type FindAllToDosUseCase struct {
	ToDoRepository repository.ToDoRepository
}

func NewFindAllToDosUseCase(todoRepository repository.ToDoRepository) *FindAllToDosUseCase {
	return &FindAllToDosUseCase{
		ToDoRepository: todoRepository,
	}
}

func (u *FindAllToDosUseCase) Execute() (*FindAllToDosOutputDTO, error) {
	res, err := u.ToDoRepository.FindAllToDos()
	if err != nil {
		return nil, err
	}
	output := make(FindAllToDosOutputDTO, len(res))
	for i, todo := range res {
		output[i] = &FindToDoOutputDTO{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
			Completed:   todo.Completed,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
		}
	}
	return &output, nil
}
