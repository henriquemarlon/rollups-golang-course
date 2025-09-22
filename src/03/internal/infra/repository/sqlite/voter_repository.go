package sqlite

import (
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
)

func (r *SQLiteRepository) CreateVoter(voter *domain.Voter) error {
	return r.db.Create(voter).Error
}

func (r *SQLiteRepository) FindVoterByID(id int) (*domain.Voter, error) {
	var voter domain.Voter
	err := r.db.First(&voter, id).Error
	if err != nil {
		return nil, err
	}
	return &voter, nil
}

func (r *SQLiteRepository) FindVoterByAddress(address Address) (*domain.Voter, error) {
	var voter domain.Voter
	err := r.db.Where("address = ?", address).First(&voter).Error
	if err != nil {
		return nil, err
	}
	return &voter, nil
}

func (r *SQLiteRepository) UpdateVoter(voter *domain.Voter) error {
	return r.db.Save(voter).Error
}

func (r *SQLiteRepository) DeleteVoter(id int) error {
	return r.db.Delete(&domain.Voter{}, id).Error
}

func (r *SQLiteRepository) HasVoted(voterID, votingID int) (bool, error) {
	var count int64
	err := r.db.Model(&domain.VotingOption{}).
		Where("voting_id = ? AND voter_id = ?", votingID, voterID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
