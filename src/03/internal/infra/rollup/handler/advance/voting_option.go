package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/usecase/voting_option"
	"github.com/rollmelette/rollmelette"
)

type VotingOptionAdvanceHandlers struct {
	VotingRepository       repository.VotingRepository
	VotingOptionRepository repository.VotingOptionRepository
}

func NewVotingOptionAdvanceHandlers(votingRepository repository.VotingRepository, votingOptionRepository repository.VotingOptionRepository) *VotingOptionAdvanceHandlers {
	return &VotingOptionAdvanceHandlers{
		VotingRepository:       votingRepository,
		VotingOptionRepository: votingOptionRepository,
	}
}

func (h *VotingOptionAdvanceHandlers) CreateVotingOption(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting_option.CreateVotingOptionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	createVotingOption := voting_option.NewCreateVotingOptionUseCase(h.VotingRepository, h.VotingOptionRepository)
	res, err := createVotingOption.Execute(ctx, &input, &metadata)
	if err != nil {
		return err
	}
	votingOptionBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting option created - "), votingOptionBytes...))
	return nil
}

func (h *VotingOptionAdvanceHandlers) DeleteVotingOption(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting_option.DeleteVotingOptionInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	deleteVotingOption := voting_option.NewDeleteVotingOptionUseCase(h.VotingOptionRepository)
	res, err := deleteVotingOption.Execute(ctx, &input, &metadata)
	if err != nil {
		return fmt.Errorf("failed to delete voting option: %w", err)
	}
	votingOptionBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting option deleted - "), votingOptionBytes...))
	return nil
}
