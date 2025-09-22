package main

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rollmelette/rollmelette"
	"github.com/stretchr/testify/suite"
)

var msgSender = common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

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

func (s *ApplicationSuite) TestVoucherDeployNFT() {
	applicationAddress := common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")

	deployNFTInput := []byte(`{"path":"deploy_nft","data":{"name":"Token","symbol":"MTK"}}`)
	newNFTOutput := s.tester.Advance(msgSender, deployNFTInput)
	s.Len(newNFTOutput.Vouchers, 1)
	s.Nil(newNFTOutput.Err)
	s.Equal(nftFactoryAddress, newNFTOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newNFT",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`

	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newNFT"].Inputs.Unpack(newNFTOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)

	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(0)), common.BytesToHash(saltBytes[:]))
}

func (s *ApplicationSuite) TestDelegatecallVoucherMintNFT() {
	uri := "https://example.com/token/1"
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	applicationAddress := common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")

	// Deploy
	deployNFTInput := []byte(`{"path":"deploy_nft","data":{"name":"Token","symbol":"TKN"}}`)
	newNFTOutput := s.tester.Advance(msgSender, deployNFTInput)
	s.Len(newNFTOutput.Vouchers, 1)
	s.Nil(newNFTOutput.Err)
	s.Equal(nftFactoryAddress, newNFTOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newNFT",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newNFT"].Inputs.Unpack(newNFTOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)

	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(0)), common.BytesToHash(saltBytes[:]))

	// Mint
	bytecode, err := getNFTBytecode()
	s.Require().NoError(err)

	stringType, _ := abi.NewType("string", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
		{Type: stringType},
		{Type: stringType},
	}.Pack(applicationAddress, "Token", "TKN")
	s.Require().NoError(err)

	nftAddress := crypto.CreateAddress2(
		nftFactoryAddress,
		common.HexToHash(strconv.Itoa(0)),
		crypto.Keccak256(append(bytecode, constructorArgs...)),
	)

	mintNFTInput := []byte(fmt.Sprintf(`{"path":"mint_nft","data":{"to":"%s","uri":"%s"}}`, to, uri))
	mintNFTOutput := s.tester.Advance(to, mintNFTInput)
	s.Nil(mintNFTOutput.Err)
	s.Len(mintNFTOutput.DelegateCallVouchers, 1)
	s.Equal(safeERC721MintAddress, mintNFTOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeMint",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"string"}
		]
	}]`
	safeMintABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err = safeMintABI.Methods["safeMint"].Inputs.Unpack(mintNFTOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(nftAddress, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(uri, unpacked[2].(string))
}

func (s *ApplicationSuite) TestDelegateCallVoucherSafeERC20Transfer() {
	amount := big.NewInt(10000)
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

	safeTransferInput := []byte(fmt.Sprintf(`{"path":"safe_erc20_transfer","data":{"token":"%s","to":"%s","amount":"%s"}}`, token, to, amount))
	safeTransferOutput := s.tester.Advance(to, safeTransferInput)
	s.Nil(safeTransferOutput.Err)
	s.Len(safeTransferOutput.DelegateCallVouchers, 1)
	s.Equal(safeERC20TransferAddress, safeTransferOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeTransfer",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"}
		]
	}]`
	safeTransferABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeTransferABI.Methods["safeTransfer"].Inputs.Unpack(safeTransferOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(amount, unpacked[2].(*big.Int))
}

func (s *ApplicationSuite) TestDelegateCallVoucherSafeERC20TransferTargeted() {
	amount := big.NewInt(10000)
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

	safeTransferTargetedInput := []byte(fmt.Sprintf(`{"path":"safe_erc20_transfer_targeted","data":{"token":"%s","to":"%s","amount":"%s"}}`, token, to, amount))
	safeTransferTargetedOutput := s.tester.Advance(to, safeTransferTargetedInput)
	s.Nil(safeTransferTargetedOutput.Err)
	s.Len(safeTransferTargetedOutput.DelegateCallVouchers, 1)
	s.Equal(safeERC20TransferAddress, safeTransferTargetedOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"safeTransferTargeted",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"address"},
			{"type":"uint256"}
		]
	}]`
	safeTransferTargetedABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := safeTransferTargetedABI.Methods["safeTransferTargeted"].Inputs.Unpack(safeTransferTargetedOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
	s.Equal(to, unpacked[2].(common.Address))
	s.Equal(amount, unpacked[3].(*big.Int))
}

func (s *ApplicationSuite) TestDelegateCallVoucherEmergencyERC20Withdraw() {
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")
	token := common.HexToAddress("0xfafafafafafafafafafafafafafafafafafafafa")

	emergencyERC20WithdrawInput := []byte(fmt.Sprintf(`{"path":"emergency_erc20_withdraw","data":{"token":"%s","to":"%s"}}`, token, to))
	emergencyERC20WithdrawOutput := s.tester.Advance(to, emergencyERC20WithdrawInput)
	s.Nil(emergencyERC20WithdrawOutput.Err)
	s.Len(emergencyERC20WithdrawOutput.DelegateCallVouchers, 1)
	s.Equal(emergencyWithdrawAddress, emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"emergencyERC20Withdraw",
		"inputs":[
			{"type":"address"},
			{"type":"address"}
		]
	}]`
	emergencyWithdrawABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := emergencyWithdrawABI.Methods["emergencyERC20Withdraw"].Inputs.Unpack(emergencyERC20WithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(token, unpacked[0].(common.Address))
	s.Equal(to, unpacked[1].(common.Address))
}

func (s *ApplicationSuite) TestDelegateCallVoucherEmergencyETHWithdraw() {
	to := common.HexToAddress("0x0000000000000000000000000000000000000001")

	emergencyETHWithdrawInput := []byte(fmt.Sprintf(`{"path":"emergency_eth_withdraw","data":{"to":"%s"}}`, to))
	emergencyETHWithdrawOutput := s.tester.Advance(to, emergencyETHWithdrawInput)
	s.Nil(emergencyETHWithdrawOutput.Err)
	s.Len(emergencyETHWithdrawOutput.DelegateCallVouchers, 1)
	s.Equal(emergencyWithdrawAddress, emergencyETHWithdrawOutput.DelegateCallVouchers[0].Destination)

	abiJSON := `[{
		"type":"function",
		"name":"emergencyETHWithdraw",
		"inputs":[
			{"type":"address"}
		]
	}]`
	emergencyETHWithdrawABI, err := abi.JSON(strings.NewReader(abiJSON))
	s.Require().NoError(err)

	unpacked, err := emergencyETHWithdrawABI.Methods["emergencyETHWithdraw"].Inputs.Unpack(emergencyETHWithdrawOutput.DelegateCallVouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(to, unpacked[0].(common.Address))
}

func (s *ApplicationSuite) TestInspectContracts() {
	applicationAddress := common.HexToAddress("0xab7528bb862fb57e8a2bcd567a2e929a0be56a5e")

	// Deploy
	deployNFTInput := []byte(`{"path":"deploy_nft","data":{"name":"Token","symbol":"TKN"}}`)
	newNFTOutput := s.tester.Advance(msgSender, deployNFTInput)
	s.Len(newNFTOutput.Vouchers, 1)
	s.Nil(newNFTOutput.Err)
	s.Equal(nftFactoryAddress, newNFTOutput.Vouchers[0].Destination)

	abiJson := `[{
		"type": "function",
		"name": "newNFT",
		"inputs": [
			{"type": "address"},
			{"type": "bytes32"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err := abiInterface.Methods["newNFT"].Inputs.Unpack(newNFTOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)

	s.Equal(applicationAddress, unpacked[0])
	saltBytes := unpacked[1].([32]byte)
	s.Equal(common.HexToHash(strconv.Itoa(0)), common.BytesToHash(saltBytes[:]))

	// Inspect
	bytecode, err := getNFTBytecode()
	s.Require().NoError(err)

	stringType, _ := abi.NewType("string", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
		{Type: stringType},
		{Type: stringType},
	}.Pack(applicationAddress, "Token", "TKN")
	s.Require().NoError(err)

	nftAddress := crypto.CreateAddress2(
		nftFactoryAddress,
		common.HexToHash(strconv.Itoa(0)),
		crypto.Keccak256(append(bytecode, constructorArgs...)),
	)

	inspectInput := []byte(`{"path":"contracts"}`)
	inspectOutput := s.tester.Inspect(inspectInput)
	s.Nil(inspectOutput.Err)
	s.Len(inspectOutput.Reports, 1)

	expectedContractsOutput := fmt.Sprintf(`[{"name":"Non Fungible Token","address":"%s"},{"name":"NFT Factory","address":"%s"},{"name":"Emergency Withdraw","address":"%s"},{"name":"Safe ERC20 Transfer","address":"%s"}]`, nftAddress, nftFactoryAddress, emergencyWithdrawAddress, safeERC20TransferAddress)
	s.Equal(expectedContractsOutput, string(inspectOutput.Reports[0].Payload))
}
