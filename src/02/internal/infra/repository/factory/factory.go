package factory

import (
	"context"
	"fmt"
	"strings"

	. "github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository/in_memory"
	"github.com/henriquemarlon/cartesi-golang-series/src/02/internal/infra/repository/sqlite"
)

func NewRepositoryFromConnectionString(ctx context.Context, conn string) (Repository, error) {
	lowerConn := strings.ToLower(conn)
	switch {
	case strings.HasPrefix(lowerConn, "memory://"):
		return newInMemoryRepository()
	case strings.HasPrefix(lowerConn, "sqlite://"):
		return newSQLiteRepository(ctx, conn)
	default:
		return nil, fmt.Errorf("unrecognized connection string format: %s", conn)
	}
}

func newInMemoryRepository() (Repository, error) {
	inMemoryRepo, err := in_memory.NewInMemoryRepository()
	if err != nil {
		return nil, err
	}

	return inMemoryRepo, nil
}

func newSQLiteRepository(ctx context.Context, conn string) (Repository, error) {
	sqliteRepo, err := sqlite.NewSQLiteRepository(ctx, conn)
	if err != nil {
		return nil, err
	}

	return sqliteRepo, nil
}
