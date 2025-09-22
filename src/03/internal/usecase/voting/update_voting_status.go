package voting

import (
	"context"
	"time"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
)

type UpdateVotingStatusUseCase struct {
	votingRepository repository.VotingRepository
}

func NewUpdateVotingStatusUseCase(votingRepository repository.VotingRepository) *UpdateVotingStatusUseCase {
	return &UpdateVotingStatusUseCase{
		votingRepository: votingRepository,
	}
}

func (u *UpdateVotingStatusUseCase) Execute(ctx context.Context) error {
	votings, err := u.votingRepository.FindAllVotings()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, voting := range votings {
		if voting.Status == domain.VotingStatusOpen && now.After(voting.EndDate) {
			voting.Status = domain.VotingStatusClosed
			if err := u.votingRepository.UpdateVoting(voting); err != nil {
				return err
			}
		}
	}

	return nil
}
