package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
)

type Application struct{}

func (a *Application) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	switch d := deposit.(type) {
	case *rollmelette.EtherDeposit:
		env.Notice([]byte(fmt.Sprintf(
			"1 - Ether balance of %s: %d, before transfer to %s",
			d.Sender.Hex(),
			env.EtherBalanceOf(d.Sender),
			common.Address{}.Hex(),
		)))
		env.EtherTransfer(d.Sender, common.Address{}, d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"2 - Balance of %s: %d before transfer to %s",
			common.Address{}.Hex(),
			env.EtherBalanceOf(common.Address{}),
			d.Sender,
		)))
		env.EtherTransfer(common.Address{}, d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"3 - Ether balance of %s: %d before withdraw",
			d.Sender.Hex(),
			env.EtherBalanceOf(d.Sender),
		)))
		env.EtherWithdraw(d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"4 - Ether balance of %s: %d after withdraw",
			d.Sender.Hex(),
			env.EtherBalanceOf(d.Sender),
		)))
		return nil

	case *rollmelette.ERC20Deposit:
		env.Notice([]byte(fmt.Sprintf(
			"1 - ERC20 balance of %s: %d before transfer to %s",
			d.Sender.Hex(),
			env.ERC20BalanceOf(d.Token, d.Sender),
			common.Address{}.Hex(),
		)))
		env.ERC20Transfer(d.Token, d.Sender, common.Address{}, d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"2 - Balance of %s: %d before transfer to %s",
			common.Address{}.Hex(),
			env.ERC20BalanceOf(d.Token, common.Address{}),
			d.Sender,
		)))
		env.ERC20Transfer(d.Token, common.Address{}, d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"3 - ERC20 balance of %s: %d before withdraw",
			d.Sender.Hex(),
			env.ERC20BalanceOf(d.Token, d.Sender),
		)))
		env.ERC20Withdraw(d.Token, d.Sender, d.Value)

		env.Notice([]byte(fmt.Sprintf(
			"4 - ERC20 balance of %s: %d after withdraw",
			d.Sender.Hex(),
			env.ERC20BalanceOf(d.Token, d.Sender),
		)))
		return nil

	default:
		env.Report([]byte(fmt.Sprintf("Unknown deposit type: %T", d)))
		return fmt.Errorf("unknown deposit type: %T", d)
	}
}

func (a *Application) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	return nil
}

func main() {
	ctx := context.Background()
	opts := rollmelette.NewRunOpts()
	app := new(Application)
	err := rollmelette.Run(ctx, opts, app)
	if err != nil {
		slog.Error("application error", "error", err)
	}
}
