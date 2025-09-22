package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
)

type FindVoterByAddressInputDTO struct {
	Address Address `json:"address" validate:"required"`
}

type FindVoterByAddressOutputDTO struct {
	Id      int     `json:"id"`
	Address Address `json:"address"`
}

type FindVoterByAddressUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewFindVoterByAddressUseCase(voterRepository repository.VoterRepository) *FindVoterByAddressUseCase {
	return &FindVoterByAddressUseCase{VoterRepository: voterRepository}
}

func (uc *FindVoterByAddressUseCase) Execute(ctx context.Context, input *FindVoterByAddressInputDTO) (*FindVoterByAddressOutputDTO, error) {
	voter, err := uc.VoterRepository.FindVoterByAddress(input.Address)
	if err != nil {
		return nil, err
	}
	return &FindVoterByAddressOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
	}, nil
}
