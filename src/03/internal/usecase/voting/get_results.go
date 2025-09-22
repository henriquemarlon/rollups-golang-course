package voting

import (
	"context"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type GetResultsInputDTO struct {
	ID int `json:"id" validate:"required"`
}

type OptionResultDTO struct {
	ID        int `json:"id"`
	VoteCount int `json:"vote_count"`
	VoterID   int `json:"voter_id"`
}

type GetResultsOutputDTO struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Status      string            `json:"status"`
	TotalVotes  int               `json:"total_votes"`
	Options     []OptionResultDTO `json:"options"`
	WinnerID    int               `json:"winner_id"`
	WinnerVotes int               `json:"winner_votes"`
}

type GetResultsUseCase struct {
	votingRepository repository.VotingRepository
}

func NewGetResultsUseCase(votingRepository repository.VotingRepository) *GetResultsUseCase {
	return &GetResultsUseCase{
		votingRepository: votingRepository,
	}
}

func (u *GetResultsUseCase) Execute(ctx context.Context, input *GetResultsInputDTO) (*GetResultsOutputDTO, error) {
	voting, err := u.votingRepository.FindVotingByID(input.ID)
	if err != nil {
		return nil, err
	}

	result := &GetResultsOutputDTO{
		ID:      voting.ID,
		Title:   voting.Title,
		Status:  string(voting.Status),
		Options: make([]OptionResultDTO, 0),
	}

	var totalVotes int
	var maxVotes int
	var winnerID int

	for _, option := range voting.Options {
		totalVotes += option.VoteCount
		if option.VoteCount > maxVotes {
			maxVotes = option.VoteCount
			winnerID = option.ID
		}

		result.Options = append(result.Options, OptionResultDTO{
			ID:        option.ID,
			VoteCount: option.VoteCount,
			VoterID:   option.VoterID,
		})
	}

	result.TotalVotes = totalVotes
	result.WinnerID = winnerID
	result.WinnerVotes = maxVotes

	return result, nil
}
