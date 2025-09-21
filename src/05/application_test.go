package main

import (
	"fmt"
	"log/slog"
	"os"
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
	bytecode, err := getNFTBytecode()
	s.Require().NoError(err)

	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
	}.Pack(applicationAddress)
	if err != nil {
		slog.Error("Failed to encode constructor args", "error", err)
		os.Exit(1)
	}

	_ = crypto.CreateAddress2(
		nftFactoryAddress,
		common.HexToHash(strconv.Itoa(0)),
		crypto.Keccak256(append(bytecode, constructorArgs...)),
	)

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

func (s *ApplicationSuite) TestVoucherMintNFT() {
	uri := "https://example.com/nft/1"
	user := common.HexToAddress("0x0000000000000000000000000000000000000001")
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

	mintNFTInput := []byte(fmt.Sprintf(`{"path":"mint_nft","data":{"to":"%s","uri":"%s"}}`, user.Hex(), uri))
	mintNFTOutput := s.tester.Advance(msgSender, mintNFTInput)
	s.Require().NoError(err)
	s.Len(mintNFTOutput.Vouchers, 1)
	s.Equal(nftAddress, mintNFTOutput.Vouchers[0].Destination)

	abiJson = `[{
		"type": "function",
		"name": "safeMint",
		"inputs": [
			{"type": "address"},
			{"type": "string"}
		]
	}]`
	abiInterface, err = abi.JSON(strings.NewReader(abiJson))
	s.Require().NoError(err)

	unpacked, err = abiInterface.Methods["safeMint"].Inputs.Unpack(mintNFTOutput.Vouchers[0].Payload[4:])
	s.Require().NoError(err)
	s.Equal(user, unpacked[0].(common.Address))
	s.Equal(uri, unpacked[1].(string))
}
