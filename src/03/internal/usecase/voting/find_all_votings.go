package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type FindAllVotingsOutputDTO struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
}

type FindAllVotingsUseCase struct {
	VotingRepository repository.VotingRepository
}

func NewFindAllVotingsUseCase(votingRepository repository.VotingRepository) *FindAllVotingsUseCase {
	return &FindAllVotingsUseCase{VotingRepository: votingRepository}
}

func (uc *FindAllVotingsUseCase) Execute(ctx context.Context) ([]*FindAllVotingsOutputDTO, error) {
	votings, err := uc.VotingRepository.FindAllVotings()
	if err != nil {
		return nil, err
	}
	var output []*FindAllVotingsOutputDTO
	for _, v := range votings {
		output = append(output, &FindAllVotingsOutputDTO{
			Id:        v.ID,
			Title:     v.Title,
			Status:    string(v.Status),
			StartDate: v.GetStartDateUnix(),
			EndDate:   v.GetEndDateUnix(),
		})
	}
	return output, nil
}
