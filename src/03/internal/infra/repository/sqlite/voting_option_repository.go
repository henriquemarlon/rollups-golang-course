package sqlite

import (
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateOption(option *domain.VotingOption) error {
	return r.db.Create(option).Error
}

func (r *SQLiteRepository) FindOptionByID(id int) (*domain.VotingOption, error) {
	var option domain.VotingOption
	err := r.db.Preload("Voting").First(&option, id).Error
	if err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *SQLiteRepository) FindAllOptionsByVotingID(votingID int) ([]*domain.VotingOption, error) {
	var options []*domain.VotingOption
	err := r.db.Preload("Voting").Where("voting_id = ?", votingID).Find(&options).Error
	if err != nil {
		return nil, err
	}
	return options, nil
}

func (r *SQLiteRepository) UpdateOption(option *domain.VotingOption) error {
	return r.db.Save(option).Error
}

func (r *SQLiteRepository) DeleteOption(id int) error {
	return r.db.Delete(&domain.VotingOption{}, id).Error
}

func (r *SQLiteRepository) IncrementVoteCount(id int, voterID int) error {
	return r.db.Model(&domain.VotingOption{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"vote_count": gorm.Expr("vote_count + ?", 1),
			"voter_id":   voterID,
		}).Error
}
