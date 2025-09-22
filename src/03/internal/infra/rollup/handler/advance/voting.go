package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/usecase/voting"
	"github.com/rollmelette/rollmelette"
)

type VotingAdvanceHandlers struct {
	CreateVotingUseCase *voting.CreateVotingUseCase
	DeleteVotingUseCase *voting.DeleteVotingUseCase
	VoteUseCase         *voting.VoteUseCase
	VotingRepository    repository.Repository
}

func NewVotingAdvanceHandlers(repo repository.Repository) *VotingAdvanceHandlers {
	return &VotingAdvanceHandlers{
		CreateVotingUseCase: voting.NewCreateVotingUseCase(repo),
		DeleteVotingUseCase: voting.NewDeleteVotingUseCase(repo),
		VoteUseCase:         voting.NewVoteUseCase(repo, repo, repo),
		VotingRepository:    repo,
	}
}

func (h *VotingAdvanceHandlers) CreateVoting(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting.CreateVotingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	res, err := h.CreateVotingUseCase.Execute(ctx, &input, &metadata)
	if err != nil {
		return fmt.Errorf("failed to create voting: %w", err)
	}
	votingBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting created - "), votingBytes...))
	return nil
}

func (h *VotingAdvanceHandlers) DeleteVoting(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting.DeleteVotingInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	res, err := h.DeleteVotingUseCase.Execute(ctx, &input, &metadata)
	if err != nil {
		return fmt.Errorf("failed to delete voting: %w", err)
	}
	votingBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voting deleted - "), votingBytes...))
	return nil
}

func (h *VotingAdvanceHandlers) Vote(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voting.VoteInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	res, err := h.VoteUseCase.Execute(input, &metadata)
	if err != nil {
		return err
	}

	voteBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("vote registered - "), voteBytes...))
	return nil
}

func (h *VotingAdvanceHandlers) UpdateStatus(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	ctx := context.Background()
	updateVotingStatus := voting.NewUpdateVotingStatusUseCase(h.VotingRepository)
	if err := updateVotingStatus.Execute(ctx); err != nil {
		return fmt.Errorf("failed to update voting status: %w", err)
	}
	env.Notice([]byte("voting status updated"))
	return nil
}
