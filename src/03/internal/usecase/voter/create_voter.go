package voter

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type CreateVoterOutputDTO struct {
	Id      int     `json:"id" validate:"required"`
	Address Address `json:"address" validate:"required"`
}

type CreateVoterUseCase struct {
	VoterRepository repository.VoterRepository
}

func NewCreateVoterUseCase(voterRepository repository.VoterRepository) *CreateVoterUseCase {
	return &CreateVoterUseCase{VoterRepository: voterRepository}
}

func (uc *CreateVoterUseCase) Execute(ctx context.Context, metadata *rollmelette.Metadata) (*CreateVoterOutputDTO, error) {
	voter, err := domain.NewVoter(Address(metadata.MsgSender))
	if err != nil {
		return nil, err
	}
	err = uc.VoterRepository.CreateVoter(voter)
	if err != nil {
		return nil, err
	}
	return &CreateVoterOutputDTO{
		Id:      voter.ID,
		Address: voter.Address,
	}, nil
}
