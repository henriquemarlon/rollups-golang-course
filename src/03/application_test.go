package main

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

var payload = common.Hex2Bytes("deadbeef")

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

type ApplicationSuite struct {
	suite.Suite
	tester *rollmelette.Tester
}

func (s *ApplicationSuite) SetupTest() {
	app := new(Application)
	s.tester = rollmelette.NewTester(app)
}

func (s *ApplicationSuite) TestERC20Deposit() {
	user := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	erc20token := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")
	depositOutput := s.tester.DepositERC20(erc20token, user, big.NewInt(10000), payload)
	s.Len(depositOutput.Notices, 4)
	s.Len(depositOutput.Vouchers, 1)
	s.Nil(depositOutput.Err)

	s.Equal(
		fmt.Sprintf(
			"1 - ERC20 balance of %s: 10000 before transfer to %s",
			user.Hex(),
			common.Address{}.Hex(),
		),
		string(depositOutput.Notices[0].Payload),
	)

	s.Equal(
		fmt.Sprintf(
			"2 - Balance of %s: 10000 before transfer to %s",
			common.Address{}.Hex(),
			user.Hex(),
		),
		string(depositOutput.Notices[1].Payload),
	)

	s.Equal(
		fmt.Sprintf(
			"3 - ERC20 balance of %s: 10000 before withdraw",
			user.Hex(),
		),
		string(depositOutput.Notices[2].Payload),
	)

	s.Equal(
		fmt.Sprintf(
			"4 - ERC20 balance of %s: 0 after withdraw",
			user.Hex(),
		),
		string(depositOutput.Notices[3].Payload),
	)

	s.Equal(erc20token, depositOutput.Vouchers[0].Destination)
	expectedWithdrawVoucherPayload := make([]byte, 0, 4+32+32)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, 0xa9, 0x05, 0x9c, 0xbb)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, make([]byte, 12)...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, user[:]...)
	expectedWithdrawVoucherPayload = append(expectedWithdrawVoucherPayload, big.NewInt(10000).FillBytes(make([]byte, 32))...)
	s.Equal(expectedWithdrawVoucherPayload, depositOutput.Vouchers[0].Payload)
	s.Equal(big.NewInt(0), depositOutput.Vouchers[0].Value)
}

func (s *ApplicationSuite) TestEtherDeposit() {
	applicationAddress := common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")

	amount := big.NewInt(10000)
	user := common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
	depositOutput := s.tester.DepositEther(user, amount, payload)
	s.Len(depositOutput.Notices, 4)
	s.Len(depositOutput.Vouchers, 1)
	s.Nil(depositOutput.Err)

	s.Equal(
		fmt.Sprintf(
			"1 - Ether balance of %s: 10000, before transfer to %s",
			user.Hex(),
			common.Address{}.Hex(),
		),
		string(depositOutput.Notices[0].Payload),
	)

	s.Equal(
		fmt.Sprintf(
			"2 - Balance of %s: 10000 before transfer to %s",
			common.Address{}.Hex(),
			user.Hex(),
		),
		string(depositOutput.Notices[1].Payload),
	)

	s.Equal(
		fmt.Sprintf(
			"3 - Ether balance of %s: 10000 before withdraw",
			user.Hex(),
		),
		string(depositOutput.Notices[2].Payload),
	)

	s.Equal(
		fmt.Sprintf(
			"4 - Ether balance of %s: 0 after withdraw",
			user.Hex(),
		),
		string(depositOutput.Notices[3].Payload),
	)

	s.Equal(applicationAddress, depositOutput.Vouchers[0].Destination)
	s.Equal(big.NewInt(10000), depositOutput.Vouchers[0].Value)
	s.Empty(depositOutput.Vouchers[0].Payload)
}
