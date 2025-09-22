package voter

import (
	"context"
	"errors"

	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type DeleteVoterInputDTO struct {
	Id int `json:"id" validate:"required"`
}

type DeleteVoterOutputDTO struct {
	Success bool `json:"success"`
}

type DeleteVoterUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewDeleteVoterUseCase(voterRepository repository.VoterRepository) *DeleteVoterUseCase {
	return &DeleteVoterUseCase{VoterRepository: voterRepository}
}

func (uc *DeleteVoterUseCase) Execute(ctx context.Context, input *DeleteVoterInputDTO, metadata *rollmelette.Metadata) (*DeleteVoterOutputDTO, error) {
	voter, err := uc.VoterRepository.FindVoterByID(input.Id)
	if err != nil {
		return &DeleteVoterOutputDTO{Success: false}, err
	}
	if voter.Address != Address(metadata.MsgSender) {
		return &DeleteVoterOutputDTO{Success: false}, errors.New("unauthorized")
	}
	err = uc.VoterRepository.DeleteVoter(input.Id)
	if err != nil {
		return &DeleteVoterOutputDTO{Success: false}, err
	}
	return &DeleteVoterOutputDTO{Success: true}, nil
}
