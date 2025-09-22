package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/cartesi/handler/advance"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/cartesi/handler/inspect"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository/factory"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/pkg/rollups"
)

var (
	infolog = log.New(os.Stderr, "[ info ] ", log.Lshortfile)
	errlog  = log.New(os.Stderr, "[ error ] ", log.Lshortfile)
)

func DappStrategy(response *rollups.FinishResponse, router *rollups.Router, ih *inspect.ToDoInspectHandlers) error {
	switch response.Type {
	case "advance_state":
		var data rollups.AdvanceResponse
		if err := json.Unmarshal(response.Data, &data); err != nil {
			return err
		}
		decodedPayload, err := hex.DecodeString(data.Payload[2:])
		if err != nil {
			return fmt.Errorf("handler: error decoding payload: %w", err)
		}
		return router.Advance(decodedPayload, data.Metadata)
	case "inspect_state":
		return ih.FindAllToDosHandler()
	}
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// toDoRepository, err := factory.NewRepositoryFromConnectionString(ctx, "memory://")
	// if err != nil {
	// 	errlog.Panicln("Failed to initialize repository", "error", err)
	// }

	toDoRepository, err := factory.NewRepositoryFromConnectionString(ctx, "sqlite:///mnt/data/database.db")
	if err != nil {
		errlog.Panicln("Failed to initialize repository", "error", err)
	}

	// Router setup and handlers registration
	defer toDoRepository.Close()

	ah := advance.NewToDoAdvanceHandlers(toDoRepository)
	if err != nil {
		errlog.Panicln("Failed to initialize advance handlers", "error", err)
	}
	infolog.Println("Advance handlers initialized")

	ih := inspect.NewToDoInspectHandlers(toDoRepository)
	if err != nil {
		errlog.Panicln("Failed to initialize inspect handlers", "error", err)
	}
	infolog.Println("Inspect handlers initialized")

	r := rollups.NewRouter()
	r.HandleAdvance("createToDo", ah.CreateToDoHandler)
	r.HandleAdvance("updateToDo", ah.UpdateToDoHandler)
	r.HandleAdvance("deleteToDo", ah.DeleteToDoHandler)
	infolog.Println("Router setup successful")

	// Polling loop ( Is there something new to process? )
	finish := rollups.FinishRequest{
		Status: "accept",
	}
	for {
		infolog.Println("Sending finish")
		res, err := rollups.SendFinish(&finish)
		if err != nil {
			errlog.Panicln("Error: error making http request: ", err)
		}
		infolog.Println("Received finish status ", strconv.Itoa(res.StatusCode))

		if res.StatusCode == 202 {
			infolog.Println("No pending rollup request, trying again")
		} else {

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				errlog.Panicln("Error: could not read response body: ", err)
			}

			var response rollups.FinishResponse
			err = json.Unmarshal(resBody, &response)
			if err != nil {
				errlog.Panicln("Error: unmarshaling body:", err)
			}
			finish.Status = "accept"

			// Strategy pattern to handle different types of requests (advance or inspect ?)
			err = DappStrategy(&response, r, ih)
			if err != nil {
				errlog.Println(err)
				finish.Status = "reject"
			}
		}
	}
}
