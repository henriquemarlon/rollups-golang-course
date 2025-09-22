package sqlite

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/domain"
)

type SQLiteRepository struct {
	db *gorm.DB
}

func (r *SQLiteRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

func NewSQLiteRepository(ctx context.Context, conn string) (*SQLiteRepository, error) {
	dbPath := strings.TrimPrefix(conn, "sqlite://")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	db = db.WithContext(ctx)

	err = db.AutoMigrate(
		&domain.Voting{},
		&domain.VotingOption{},
		&domain.Voter{},
	)
	if err != nil {
		return nil, err
	}

	return &SQLiteRepository{db: db}, nil
}
