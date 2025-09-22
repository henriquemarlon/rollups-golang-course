package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
)

type FindVoterByIDInputDTO struct {
	Id int `json:"id" validate:"required"`
}

type FindVoterByIDOutputDTO struct {
	Id      int     `json:"id"`
	Address Address `json:"address"`
}

type FindVoterByIDUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewFindVoterByIDUseCase(voterRepository repository.VoterRepository) *FindVoterByIDUseCase {
	return &FindVoterByIDUseCase{VoterRepository: voterRepository}
}

func (uc *FindVoterByIDUseCase) Execute(ctx context.Context, input *FindVoterByIDInputDTO) (*FindVoterByIDOutputDTO, error) {
	voter, err := uc.VoterRepository.FindVoterByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindVoterByIDOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
	}, nil
}
