package voting

import (
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
	"github.com/rollmelette/rollmelette"
)

type VoteInputDTO struct {
	VotingID int `json:"voting_id" validate:"required"`
	OptionID int `json:"option_id" validate:"required"`
}

type VoteOutputDTO struct {
	VotingID  int     `json:"voting_id"`
	OptionID  int     `json:"option_id"`
	Voter     Address `json:"voter"`
	VoteCount int     `json:"vote_count"`
}

type VoteUseCase struct {
	VotingRepository       repository.VotingRepository
	VotingOptionRepository repository.VotingOptionRepository
	VoterRepository        repository.VoterRepository
}

func NewVoteUseCase(
	votingRepository repository.VotingRepository,
	votingOptionRepository repository.VotingOptionRepository,
	voterRepository repository.VoterRepository,
) *VoteUseCase {
	return &VoteUseCase{
		VotingRepository:       votingRepository,
		VotingOptionRepository: votingOptionRepository,
		VoterRepository:        voterRepository,
	}
}

func (u *VoteUseCase) Execute(input VoteInputDTO, metadata *rollmelette.Metadata) (*VoteOutputDTO, error) {
	voting, err := u.VotingRepository.FindVotingByID(input.VotingID)
	if err != nil {
		return nil, fmt.Errorf("failed to find voting: %w", err)
	}

	if voting.Status != domain.VotingStatusOpen {
		return nil, domain.ErrVotingClosed
	}

	voter, err := u.VoterRepository.FindVoterByAddress(Address(metadata.MsgSender))
	if err != nil {
		return nil, domain.ErrVoterNotFound
	}

	hasVoted, err := u.VoterRepository.HasVoted(voter.ID, voting.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if voter has voted: %w", err)
	}
	if hasVoted {
		return nil, domain.ErrAlreadyVoted
	}

	option, err := u.VotingOptionRepository.FindOptionByID(input.OptionID)
	if err != nil {
		return nil, domain.ErrOptionNotFound
	}

	if option.VotingID != voting.ID {
		return nil, domain.ErrInvalidOption
	}

	err = u.VotingOptionRepository.IncrementVoteCount(option.ID, voter.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to increment vote count: %w", err)
	}

	return &VoteOutputDTO{
		VotingID:  voting.ID,
		OptionID:  option.ID,
		Voter:     voter.Address,
		VoteCount: option.VoteCount + 1,
	}, nil
}
