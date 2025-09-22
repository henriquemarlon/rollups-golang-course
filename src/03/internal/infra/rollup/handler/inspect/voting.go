package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/usecase/voting"
	"github.com/rollmelette/rollmelette"
)

type VotingInspectHandlers struct {
	VotingRepository       repository.VotingRepository
	VotingOptionRepository repository.VotingOptionRepository
}

func NewVotingInspectHandlers(votingRepository repository.VotingRepository, votingOptionRepository repository.VotingOptionRepository) *VotingInspectHandlers {
	return &VotingInspectHandlers{
		VotingRepository:       votingRepository,
		VotingOptionRepository: votingOptionRepository,
	}
}

func (h *VotingInspectHandlers) FindAllVotings(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllVotings := voting.NewFindAllVotingsUseCase(h.VotingRepository)
	votings, err := findAllVotings.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all votings: %w", err)
	}
	votingsBytes, err := json.Marshal(votings)
	if err != nil {
		return fmt.Errorf("failed to marshal votings: %w", err)
	}
	env.Report(votingsBytes)
	return nil
}

func (h *VotingInspectHandlers) FindVotingByID(env rollmelette.EnvInspector, payload []byte) error {
	var input voting.FindVotingByIDInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findVotingByID := voting.NewFindVotingByIDUseCase(h.VotingRepository)
	votingRes, err := findVotingByID.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find voting by id: %w", err)
	}
	votingBytes, err := json.Marshal(votingRes)
	if err != nil {
		return fmt.Errorf("failed to marshal voting: %w", err)
	}
	env.Report(votingBytes)
	return nil
}

func (h *VotingInspectHandlers) FindAllActiveVotings(env rollmelette.EnvInspector, payload []byte) error {
	ctx := context.Background()
	findAllActiveVotings := voting.NewFindAllActiveVotingsUseCase(h.VotingRepository)
	votings, err := findAllActiveVotings.Execute(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all active votings: %w", err)
	}
	votingsBytes, err := json.Marshal(votings)
	if err != nil {
		return fmt.Errorf("failed to marshal votings: %w", err)
	}
	env.Report(votingsBytes)
	return nil
}

func (h *VotingInspectHandlers) GetVotingResults(env rollmelette.EnvInspector, payload []byte) error {
	var input voting.GetVotingResultsInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	getResults := voting.NewGetVotingResultsUseCase(h.VotingOptionRepository)
	results, err := getResults.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to get voting results: %w", err)
	}
	resultsBytes, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}
	env.Report(resultsBytes)
	return nil
}

func (h *VotingInspectHandlers) GetResults(env rollmelette.EnvInspector, payload []byte) error {
	var input voting.GetResultsInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	getResults := voting.NewGetResultsUseCase(h.VotingRepository)
	result, err := getResults.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to get voting results: %w", err)
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	env.Report(resultBytes)
	return nil
}
