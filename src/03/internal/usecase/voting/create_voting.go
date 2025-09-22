package voting

import (
	"context"
	"time"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type CreateVotingInputDTO struct {
	Title     string `json:"title" validate:"required"`
	StartDate int64  `json:"start_date" validate:"required"`
	EndDate   int64  `json:"end_date" validate:"required"`
}

type CreateVotingOutputDTO struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Creator   Address `json:"creator"`
	Status    string  `json:"status"`
	StartDate int64   `json:"start_date"`
	EndDate   int64   `json:"end_date"`
}

type CreateVotingUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewCreateVotingUseCase(votingRepository repository.VotingRepository) *CreateVotingUseCase {
	return &CreateVotingUseCase{VotingRepository: votingRepository}
}

func (uc *CreateVotingUseCase) Execute(ctx context.Context, input *CreateVotingInputDTO, metadata *rollmelette.Metadata) (*CreateVotingOutputDTO, error) {
	startDate := time.Unix(input.StartDate, 0)
	endDate := time.Unix(input.EndDate, 0)
	voting, err := domain.NewVoting(input.Title, Address(metadata.MsgSender), startDate, endDate)
	if err != nil {
		return nil, err
	}
	err = uc.VotingRepository.CreateVoting(voting)
	if err != nil {
		return nil, err
	}
	return &CreateVotingOutputDTO{
		Id:        voting.ID,
		Title:     voting.Title,
		Creator:   voting.Creator,
		Status:    string(voting.Status),
		StartDate: voting.GetStartDateUnix(),
		EndDate:   voting.GetEndDateUnix(),
	}, nil
}
