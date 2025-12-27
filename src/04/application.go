package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/rollmelette/rollmelette"
)

var (
	nftAddress        common.Address
	nftFactoryAddress = common.HexToAddress("0x24D451CC632BE1FF86f0AaEaAC026261fFd889A0") // NOTE: this address is computed from the salt "1596"
)

type Application struct{}

func (a *Application) Advance(
	env rollmelette.Env,
	metadata rollmelette.Metadata,
	deposit rollmelette.Deposit,
	payload []byte,
) error {
	var input struct {
		Path string          `json:"path" validate:"required"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}

	validator := validator.New()
	if err := validator.Struct(input); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	switch input.Path {
	case "deploy_nft":
		var data struct {
			Name   string `json:"name" validate:"required"`
			Symbol string `json:"symbol" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		deployNFTPayload, err := buildDeployNFTVoucherPayload(
			metadata.AppContract,
			common.HexToHash(strconv.Itoa(metadata.Index)),
			data.Name,
			data.Symbol,
		)
		if err != nil {
			return err
		}
		env.Voucher(nftFactoryAddress, big.NewInt(0), deployNFTPayload)
		return nil

	case "mint_nft":
		var data struct {
			To  common.Address `json:"to" validate:"required"`
			URI string         `json:"uri" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		if err := validator.Struct(data); err != nil {
			return fmt.Errorf("failed to validate input: %w", err)
		}
		voucher, err := buildMintNFTVoucherPayload(data.To, data.URI)
		if err != nil {
			return err
		}
		env.Voucher(nftAddress, big.NewInt(0), voucher)
		return nil

	default:
		env.Report([]byte(fmt.Sprintf("Unknown path: %s", input.Path)))
		return fmt.Errorf("unknown path: %s", input.Path)
	}
}

func (a *Application) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	var data struct {
		Path string `json:"path" validate:"required"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	validator := validator.New()
	if err := validator.Struct(data); err != nil {
		return fmt.Errorf("failed to validate input: %w", err)
	}

	switch data.Path {
	case "contracts":
		contractsJson := fmt.Sprintf(
			`[{"name":"Non Fungible Token","address":"%s"},{"name":"NFT Factory","address":"%s"}]`,
			nftAddress,
			nftFactoryAddress,
		)
		env.Report([]byte(contractsJson))
		return nil

	default:
		env.Report([]byte(fmt.Sprintf("Unknown path: %s", data.Path)))
		return fmt.Errorf("unknown path: %s", data.Path)
	}
}

func buildDeployNFTVoucherPayload(initialOwner common.Address, salt common.Hash, name string, symbol string) ([]byte, error) {
	bytecode, err := getNFTBytecode()
	if err != nil {
		return nil, err
	}
	stringType, _ := abi.NewType("string", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	constructorArgs, err := abi.Arguments{
		{Type: addressType},
		{Type: stringType},
		{Type: stringType},
	}.Pack(initialOwner, name, symbol)
	if err != nil {
		return nil, fmt.Errorf("error encoding constructor args: %w", err)
	}
	nftAddress = crypto.CreateAddress2(
		nftFactoryAddress,
		salt,
		crypto.Keccak256(append(bytecode, constructorArgs...)),
	)

	abiJSON := `[{
		"type":"function",
		"name":"newNFT",
		"inputs":[
			{"type":"address"},
			{"type":"bytes32"},
			{"type":"string"},
			{"type":"string"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	voucher, err := abiInterface.Pack(
		"newNFT",
		initialOwner,
		salt,
		name,
		symbol,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %w", err)
	}
	return voucher, nil
}

func buildMintNFTVoucherPayload(to common.Address, uri string) ([]byte, error) {
	abiJSON := `[{
		"type":"function",
		"name":"safeMint",
		"inputs":[
			{"type":"address"},
			{"type":"string"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	voucher, err := abiInterface.Pack("safeMint", to, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %w", err)
	}
	return voucher, nil
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
