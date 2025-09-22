package inspect

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/usecase/voting_option"
	"github.com/rollmelette/rollmelette"
)

type VotingOptionInspectHandlers struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewVotingOptionInspectHandlers(votingOptionRepository repository.VotingOptionRepository) *VotingOptionInspectHandlers {
	return &VotingOptionInspectHandlers{
		VotingOptionRepository: votingOptionRepository,
	}
}

func (h *VotingOptionInspectHandlers) FindVotingOptionByID(env rollmelette.EnvInspector, payload []byte) error {
	var input voting_option.FindVotingOptionByIDInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findVotingOptionByID := voting_option.NewFindVotingOptionByIDUseCase(h.VotingOptionRepository)
	votingOptionRes, err := findVotingOptionByID.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find voting option by id: %w", err)
	}
	votingOptionBytes, err := json.Marshal(votingOptionRes)
	if err != nil {
		return fmt.Errorf("failed to marshal voting option: %w", err)
	}
	env.Report(votingOptionBytes)
	return nil
}

func (h *VotingOptionInspectHandlers) FindAllOptionsByVotingID(env rollmelette.EnvInspector, payload []byte) error {
	var input voting_option.FindAllOptionsByVotingIDInputDTO
	if err := json.Unmarshal(payload, &input); err != nil {
		return fmt.Errorf("failed to unmarshal input: %w", err)
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	ctx := context.Background()
	findAllOptionsByVotingID := voting_option.NewFindAllOptionsByVotingIDUseCase(h.VotingOptionRepository)
	options, err := findAllOptionsByVotingID.Execute(ctx, &input)
	if err != nil {
		return fmt.Errorf("failed to find all options by voting id: %w", err)
	}
	optionsBytes, err := json.Marshal(options)
	if err != nil {
		return fmt.Errorf("failed to marshal options: %w", err)
	}
	env.Report(optionsBytes)
	return nil
}
