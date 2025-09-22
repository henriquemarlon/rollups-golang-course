package voting_option

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type FindVotingOptionByIDInputDTO struct {
	Id int `json:"id" validate:"required"`
}

type FindVotingOptionByIDOutputDTO struct {
	Id        int `json:"id"`
	VotingID  int `json:"voting_id"`
	VoteCount int `json:"vote_count"`
}

type FindVotingOptionByIDUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewFindVotingOptionByIDUseCase(votingOptionRepository repository.VotingOptionRepository) *FindVotingOptionByIDUseCase {
	return &FindVotingOptionByIDUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *FindVotingOptionByIDUseCase) Execute(ctx context.Context, input *FindVotingOptionByIDInputDTO) (*FindVotingOptionByIDOutputDTO, error) {
	option, err := uc.VotingOptionRepository.FindOptionByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindVotingOptionByIDOutputDTO{
		Id:        option.ID,
		VotingID:  option.VotingID,
		VoteCount: option.VoteCount,
	}, nil
}
