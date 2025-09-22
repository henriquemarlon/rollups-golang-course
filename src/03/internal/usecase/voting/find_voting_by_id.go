package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type FindVotingByIDInputDTO struct {
	Id int `json:"id" validate:"required"`
}

type FindVotingByIDOutputDTO struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
}

type FindVotingByIDUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewFindVotingByIDUseCase(votingRepository repository.VotingRepository) *FindVotingByIDUseCase {
	return &FindVotingByIDUseCase{VotingRepository: votingRepository}
}

func (uc *FindVotingByIDUseCase) Execute(ctx context.Context, input *FindVotingByIDInputDTO) (*FindVotingByIDOutputDTO, error) {
	voting, err := uc.VotingRepository.FindVotingByID(input.Id)
	if err != nil {
		return nil, err
	}
	return &FindVotingByIDOutputDTO{
		Id:        voting.ID,
		Title:     voting.Title,
		Status:    string(voting.Status),
		StartDate: voting.GetStartDateUnix(),
		EndDate:   voting.GetEndDateUnix(),
	}, nil
}
