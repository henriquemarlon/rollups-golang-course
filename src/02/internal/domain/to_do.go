package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidToDo = errors.New("invalid todo")
	ErrNotFound    = errors.New("todo not found")
)

type ToDo struct {
	Id          uint   `json:"id" gorm:"primaryKey"`
	Title       string `json:"title" gorm:"type:text;not null"`
	Description string `json:"description" gorm:"type:text;not null"`
	Completed   bool   `json:"completed" gorm:"default:false"`
	CreatedAt   uint64 `json:"created_at,omitempty" gorm:"not null"`
	UpdatedAt   uint64 `json:"updated_at,omitempty" gorm:"default:0"`
}

func NewToDo(title string, description string, createdAt uint64) (*ToDo, error) {
	toDo := &ToDo{
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
	}
	if err := toDo.Validate(); err != nil {
		return nil, err
	}
	return toDo, nil
}

func (t *ToDo) Validate() error {
	if t.Title == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrInvalidToDo)
	}
	if t.Description == "" {
		return fmt.Errorf("%w: description cannot be empty", ErrInvalidToDo)
	}
	return nil
}
