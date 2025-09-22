package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type GetVotingResultsInputDTO struct {
	VotingID int `json:"voting_id" validate:"required"`
}

type VotingOptionResultDTO struct {
	OptionID  int `json:"option_id"`
	VoteCount int `json:"vote_count"`
}

type GetVotingResultsOutputDTO struct {
	VotingID int                     `json:"voting_id"`
	Results  []VotingOptionResultDTO `json:"results"`
}

type GetVotingResultsUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewGetVotingResultsUseCase(optionRepo repository.VotingOptionRepository) *GetVotingResultsUseCase {
	return &GetVotingResultsUseCase{VotingOptionRepository: optionRepo}
}

func (uc *GetVotingResultsUseCase) Execute(ctx context.Context, input *GetVotingResultsInputDTO) (*GetVotingResultsOutputDTO, error) {
	options, err := uc.VotingOptionRepository.FindAllOptionsByVotingID(input.VotingID)
	if err != nil {
		return nil, err
	}
	results := make([]VotingOptionResultDTO, 0, len(options))
	for _, opt := range options {
		results = append(results, VotingOptionResultDTO{
			OptionID:  opt.ID,
			VoteCount: opt.VoteCount,
		})
	}
	return &GetVotingResultsOutputDTO{
		VotingID: input.VotingID,
		Results:  results,
	}, nil
}
