package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type FindAllOptionsByVotingIDInputDTO struct {
	VotingID int `json:"voting_id" validate:"required"`
}

type FindAllOptionsByVotingIDOutputDTO struct {
	Id        int `json:"id"`
	VotingID  int `json:"voting_id"`
	VoteCount int `json:"vote_count"`
}

type FindAllOptionsByVotingIDUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewFindAllOptionsByVotingIDUseCase(votingOptionRepository repository.VotingOptionRepository) *FindAllOptionsByVotingIDUseCase {
	return &FindAllOptionsByVotingIDUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *FindAllOptionsByVotingIDUseCase) Execute(ctx context.Context, input *FindAllOptionsByVotingIDInputDTO) ([]*FindAllOptionsByVotingIDOutputDTO, error) {
	options, err := uc.VotingOptionRepository.FindAllOptionsByVotingID(input.VotingID)
	if err != nil {
		return nil, err
	}
	var output []*FindAllOptionsByVotingIDOutputDTO
	for _, o := range options {
		output = append(output, &FindAllOptionsByVotingIDOutputDTO{
			Id:        o.ID,
			VotingID:  o.VotingID,
			VoteCount: o.VoteCount,
		})
	}
	return output, nil
}
