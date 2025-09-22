package domain

import (
	"errors"
	"fmt"

	. "github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/custom_type"
)

var (
	ErrInvalidVoter  = errors.New("invalid voter")
	ErrVoterNotFound = errors.New("voter not found")
)

type Voter struct {
	ID      int     `gorm:"primaryKey;autoIncrement"`
	Address Address `gorm:"not null;uniqueIndex"`
}

func NewVoter(address Address) (*Voter, error) {
	voter := &Voter{
		Address: address,
	}
	if err := voter.validate(); err != nil {
		return nil, err
	}
	return voter, nil
}

func (v *Voter) validate() error {
	if v.Address == (Address{}) {
		return fmt.Errorf("%w: address cannot be empty", ErrInvalidVoter)
	}
	return nil
}
