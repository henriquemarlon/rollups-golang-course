package rollup

import (
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/repository"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/rollup/handler/advance"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/internal/infra/rollup/handler/inspect"
	"github.com/henriquemarlon/cartesi-golang-series/src/03/pkg/router"
)

type CreateInfo struct {
	Repo repository.Repository
}

func Create(c *CreateInfo) *router.Router {
	votingAdvanceHandlers := advance.NewVotingAdvanceHandlers(c.Repo)
	votingInspectHandlers := inspect.NewVotingInspectHandlers(c.Repo, c.Repo)

	voterAdvanceHandlers := advance.NewVoterAdvanceHandlers(c.Repo)
	voterInspectHandlers := inspect.NewVoterInspectHandlers(c.Repo)

	votingOptionAdvanceHandlers := advance.NewVotingOptionAdvanceHandlers(c.Repo, c.Repo)
	votingOptionInspectHandlers := inspect.NewVotingOptionInspectHandlers(c.Repo)

	r := router.NewRouter()
	r.Use(router.LoggingMiddleware)
	r.Use(router.ErrorHandlingMiddleware)

	{
		votingGroup := r.Group("voting")
		votingGroup.HandleAdvance("create", votingAdvanceHandlers.CreateVoting)
		votingGroup.HandleAdvance("delete", votingAdvanceHandlers.DeleteVoting)
		votingGroup.HandleAdvance("vote", votingAdvanceHandlers.Vote)
		votingGroup.HandleAdvance("update-status", votingAdvanceHandlers.UpdateStatus)

		votingGroup.HandleInspect("", votingInspectHandlers.FindAllVotings)
		votingGroup.HandleInspect("id", votingInspectHandlers.FindVotingByID)
		votingGroup.HandleInspect("active", votingInspectHandlers.FindAllActiveVotings)
		votingGroup.HandleInspect("results", votingInspectHandlers.GetResults)
	}
	
	{
		voterGroup := r.Group("voter")
		voterGroup.HandleAdvance("create", voterAdvanceHandlers.CreateVoter)
		voterGroup.HandleAdvance("delete", voterAdvanceHandlers.DeleteVoter)

		voterGroup.HandleInspect("id", voterInspectHandlers.FindVoterByID)
		voterGroup.HandleInspect("address", voterInspectHandlers.FindVoterByAddress)
	}

	{
		votingOptionGroup := r.Group("voting-option")
		votingOptionGroup.HandleAdvance("create", votingOptionAdvanceHandlers.CreateVotingOption)
		votingOptionGroup.HandleAdvance("delete", votingOptionAdvanceHandlers.DeleteVotingOption)

		votingOptionGroup.HandleInspect("id", votingOptionInspectHandlers.FindVotingOptionByID)
		votingOptionGroup.HandleInspect("voting", votingOptionInspectHandlers.FindAllOptionsByVotingID)
	}
	return r
}
