package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type FindAllActiveVotingsOutputDTO struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
}

type FindAllActiveVotingsUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewFindAllActiveVotingsUseCase(votingRepository repository.VotingRepository) *FindAllActiveVotingsUseCase {
	return &FindAllActiveVotingsUseCase{VotingRepository: votingRepository}
}

func (uc *FindAllActiveVotingsUseCase) Execute(ctx context.Context) ([]*FindAllActiveVotingsOutputDTO, error) {
	votings, err := uc.VotingRepository.FindAllActiveVotings()
	if err != nil {
		return nil, err
	}
	var output []*FindAllActiveVotingsOutputDTO
	for _, v := range votings {
		output = append(output, &FindAllActiveVotingsOutputDTO{
			Id:        v.ID,
			Title:     v.Title,
			Status:    string(v.Status),
			StartDate: v.GetStartDateUnix(),
			EndDate:   v.GetEndDateUnix(),
		})
	}
	return output, nil
}
