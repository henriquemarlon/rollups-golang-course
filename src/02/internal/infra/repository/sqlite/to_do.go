package sqlite

import (
	"fmt"

	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/domain"
	"gorm.io/gorm"
)

func (r *SQLiteRepository) CreateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	if err := r.Db.Create(input).Error; err != nil {
		return nil, fmt.Errorf("failed to create src/02: %w", err)
	}
	return input, nil
}

func (r *SQLiteRepository) FindAllToDos() ([]*domain.ToDo, error) {
	var toDos []*domain.ToDo
	if err := r.Db.Find(&toDos).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find all src/02s: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find all src/02s: %w", err)
	}
	return toDos, nil
}

func (r *SQLiteRepository) UpdateToDo(input *domain.ToDo) (*domain.ToDo, error) {
	if err := r.Db.Updates(&input).Error; err != nil {
		return nil, fmt.Errorf("failed to update src/02: %w", err)
	}
	toDo, err := r.findToDoById(input.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to update src/02: %w", err)
	}
	return toDo, nil
}

func (r *SQLiteRepository) DeleteToDo(id uint) error {
	if err := r.Db.Delete(&domain.ToDo{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to delete src/02: %w", domain.ErrNotFound)
		}
		return fmt.Errorf("failed to delete src/02: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) findToDoById(id uint) (*domain.ToDo, error) {
	var toDo domain.ToDo
	if err := r.Db.First(&toDo, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to find src/02 by id: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to find src/02 by id: %w", err)
	}
	return &toDo, nil
}
