package advance

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/usecase/voter"
	"github.com/rollmelette/rollmelette"
)

type VoterAdvanceHandlers struct {
	VoterRepository repository.VoterRepository
}

func NewVoterAdvanceHandlers(voterRepository repository.VoterRepository) *VoterAdvanceHandlers {
	return &VoterAdvanceHandlers{
		VoterRepository: voterRepository,
	}
}

func (h *VoterAdvanceHandlers) CreateVoter(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	ctx := context.Background()
	createVoter := voter.NewCreateVoterUseCase(h.VoterRepository)
	res, err := createVoter.Execute(ctx, &metadata)
	if err != nil {
		return fmt.Errorf("failed to create voter: %w", err)
	}
	voterBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voter created - "), voterBytes...))
	return nil
}

func (h *VoterAdvanceHandlers) DeleteVoter(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	var input voter.DeleteVoterInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	deleteVoter := voter.NewDeleteVoterUseCase(h.VoterRepository)
	res, err := deleteVoter.Execute(ctx, &input, &metadata)
	if err != nil {
		return fmt.Errorf("failed to delete voter: %w", err)
	}
	voterBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	env.Notice(append([]byte("voter deleted - "), voterBytes...))
	return nil
}
