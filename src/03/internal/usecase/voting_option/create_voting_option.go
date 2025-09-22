package voting_option

import (
	"context"
	"errors"

	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type CreateVotingOptionInputDTO struct {
	VotingID int `json:"voting_id" validate:"required"`
}

type CreateVotingOptionOutputDTO struct {
	Id       int `json:"id"`
	VotingID int `json:"voting_id"`
}

type CreateVotingOptionUseCase struct {
	VotingRepository       repository.VotingRepository
	VotingOptionRepository repository.VotingOptionRepository
}

func NewCreateVotingOptionUseCase(votingRepository repository.VotingRepository, votingOptionRepository repository.VotingOptionRepository) *CreateVotingOptionUseCase {
	return &CreateVotingOptionUseCase{
		VotingRepository:       votingRepository,
		VotingOptionRepository: votingOptionRepository,
	}
}

func (uc *CreateVotingOptionUseCase) Execute(ctx context.Context, input *CreateVotingOptionInputDTO, metadata *rollmelette.Metadata) (*CreateVotingOptionOutputDTO, error) {
	if input.VotingID <= 0 {
		return nil, domain.ErrInvalidVotingOption
	}

	voting, err := uc.VotingRepository.FindVotingByID(input.VotingID)
	if err != nil {
		return nil, err
	}
	if voting.Creator != Address(metadata.MsgSender) {
		return nil, errors.New("unauthorized")
	}
	option, err := domain.NewVotingOption(input.VotingID)
	if err != nil {
		return nil, err
	}
	err = uc.VotingOptionRepository.CreateOption(option)
	if err != nil {
		return nil, err
	}
	return &CreateVotingOptionOutputDTO{
		Id:       option.ID,
		VotingID: option.VotingID,
	}, nil
}
