package root

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository/factory"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/rollup"
	"github.com/rollmelette/rollmelette"
	"github.com/spf13/cobra"
)

const (
	CMD_NAME = "rollup"
)

var (
	useMemoryDB bool
	Cmd         = &cobra.Command{
		Use:   "voting-" + CMD_NAME,
		Short: "Runs Voting Rollup",
		Long:  `Cartesi Rollup Application for voting`,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	repo, err := factory.NewRepositoryFromConnectionString(ctx, "sqlite:///mnt/data/voting.db")
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database initialized")

	defer repo.Close()

	createInfo := &rollup.CreateInfo{
		Repo: repo,
	}

	r := rollup.Create(createInfo)
	opts := rollmelette.NewRunOpts()
	if err := rollmelette.Run(ctx, opts, r); err != nil {
		slog.Error("Failed to run rollmelette", "error", err)
		os.Exit(1)
	}
}
