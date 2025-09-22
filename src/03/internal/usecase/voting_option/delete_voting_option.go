package voting_option

import (
	"context"
	"errors"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type DeleteVotingOptionInputDTO struct {
	Id int `json:"id" validate:"required"`
}

type DeleteVotingOptionOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVotingOptionUseCase struct {
	VotingOptionRepository repository.VotingOptionRepository
}

func NewDeleteVotingOptionUseCase(votingOptionRepository repository.VotingOptionRepository) *DeleteVotingOptionUseCase {
	return &DeleteVotingOptionUseCase{VotingOptionRepository: votingOptionRepository}
}

func (uc *DeleteVotingOptionUseCase) Execute(ctx context.Context, input *DeleteVotingOptionInputDTO, metadata *rollmelette.Metadata) (*DeleteVotingOptionOutputDTO, error) {
	votingOption, err := uc.VotingOptionRepository.FindOptionByID(input.Id)
	if err != nil {
		return &DeleteVotingOptionOutputDTO{Success: false}, err
	}
	if votingOption.Voting == nil {
		return &DeleteVotingOptionOutputDTO{Success: false}, errors.New("voting not found")
	}
	if votingOption.Voting.Creator != Address(metadata.MsgSender) {
		return &DeleteVotingOptionOutputDTO{Success: false}, errors.New("unauthorized")
	}
	err = uc.VotingOptionRepository.DeleteOption(input.Id)
	if err != nil {
		return &DeleteVotingOptionOutputDTO{Success: false}, err
	}
	return &DeleteVotingOptionOutputDTO{Success: true}, nil
}
