package voting

import (
	"context"
	"errors"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type DeleteVotingInputDTO struct {
	Id int `json:"id" validate:"required"`
}

type DeleteVotingOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewDeleteVotingUseCase(votingRepository repository.VotingRepository) *DeleteVotingUseCase {
	return &DeleteVotingUseCase{VotingRepository: votingRepository}
}

func (uc *DeleteVotingUseCase) Execute(ctx context.Context, input *DeleteVotingInputDTO, metadata *rollmelette.Metadata) (*DeleteVotingOutputDTO, error) {
	voting, err := uc.VotingRepository.FindVotingByID(input.Id)
	if err != nil {
		return &DeleteVotingOutputDTO{Success: false}, err
	}
	if voting.Creator != Address(metadata.MsgSender) {
		return &DeleteVotingOutputDTO{Success: false}, errors.New("unauthorized")
	}
	err = uc.VotingRepository.DeleteVoting(input.Id)
	if err != nil {
		return &DeleteVotingOutputDTO{Success: false}, err
	}
	return &DeleteVotingOutputDTO{Success: true}, nil
}
