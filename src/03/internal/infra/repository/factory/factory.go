package factory

import (
	"context"
	"fmt"
	"strings"

	. "github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository/sqlite"
)

func NewRepositoryFromConnectionString(ctx context.Context, conn string) (Repository, error) {
	lowerConn := strings.ToLower(conn)
	switch {
	case strings.HasPrefix(lowerConn, "sqlite://"):
		return newSQLiteRepository(ctx, conn)
	default:
		return nil, fmt.Errorf("unrecognized connection string format: %s", conn)
	}
}

func newSQLiteRepository(ctx context.Context, conn string) (Repository, error) {
	sqliteRepo, err := sqlite.NewSQLiteRepository(ctx, conn)
	if err != nil {
		return nil, err
	}

	return sqliteRepo, nil
}
