package inspect

import (
	"encoding/json"

	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/usecase"
	rollups "github.com/henriquemarlon/cartesi-golang-series/src/02/pkg/rollups"
)

type ToDoInspectHandlers struct {
	ToDoRepository repository.ToDoRepository
}

func NewToDoInspectHandlers(toDoRepository repository.ToDoRepository) *ToDoInspectHandlers {
	return &ToDoInspectHandlers{
		ToDoRepository: toDoRepository,
	}
}

func (h *ToDoInspectHandlers) FindAllToDosHandler() error {
	findAllToDos := usecase.NewFindAllToDosUseCase(h.ToDoRepository)
	res, err := findAllToDos.Execute()
	if err != nil {
		return err
	}
	toDos, err := json.Marshal(res)
	if err != nil {
		return err
	}
	rollups.SendReport(&rollups.ReportRequest{
		Payload: rollups.Str2Hex(string(toDos)),
	})
	return nil
}
