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
	"github.com/holiman/uint256"
	"github.com/rollmelette/rollmelette"
)

var (
	nftAddress               common.Address
	nftFactoryAddress        = common.HexToAddress("0x24D451CC632BE1FF86f0AaEaAC026261fFd889A0") // NOTE: this address is computed from the salt "1596"
	safeERC721MintAddress    = common.HexToAddress("0x4F85347240488E62ab1C6169Cbc532A09223efa4") // NOTE: this address is computed from the salt "1596"
	safeERC20TransferAddress = common.HexToAddress("0x86E244fbb3243f19492A3d61336e285bbf8E6154") // NOTE: this address is computed from the salt "1596"
	emergencyWithdrawAddress = common.HexToAddress("0xA716b0bE3a59b05A307b98c6bAf9d21dF796F37d") // NOTE: this address is computed from the salt "1596"
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
			Name  string `json:"name" validate:"required"`
			Symbol string `json:"symbol" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		bytecode, err := getNFTBytecode()
		if err != nil {
			return err
		}
		stringType, _ := abi.NewType("string", "", nil)
		addressType, _ := abi.NewType("address", "", nil)
		constructorArgs, err := abi.Arguments{
			{Type: addressType},
			{Type: stringType},
			{Type: stringType},
		}.Pack(metadata.AppContract, data.Name, data.Symbol)
		if err != nil {
			return fmt.Errorf("error encoding constructor args: %w", err)
		}
		nftAddress = crypto.CreateAddress2(
			nftFactoryAddress,
			common.HexToHash(strconv.Itoa(metadata.Index)),
			crypto.Keccak256(append(bytecode, constructorArgs...)),
		)
		deployNFTPayload, err := deployNFT(
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
		mintNFTPayload, err := mintNFT(nftAddress, data.To, data.URI)
		if err != nil {
			return err
		}
		env.DelegateCallVoucher(safeERC721MintAddress, mintNFTPayload)
		return nil

	case "safe_transfer":
		var data struct {
			Token  common.Address `json:"token" validate:"required"`
			To     common.Address `json:"to" validate:"required"`
			Amount *uint256.Int   `json:"amount" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		if err := validator.Struct(data); err != nil {
			return fmt.Errorf("failed to validate input: %w", err)
		}
		abiJSON := `[{
			"type":"function",
			"name":"safeTransfer",
			"inputs":[
				{"type":"address"},
				{"type":"address"},
				{"type":"uint256"}
			]
		}]`
		abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return err
		}
		delegateCallPayload, err := abiInterface.Pack(
			"safeTransfer",
			data.Token,
			data.To,
			data.Amount.ToBig(),
		)
		if err != nil {
			return err
		}
		env.SetERC20Balance(
			data.Token,
			metadata.MsgSender,
			new(big.Int).Sub(
				env.ERC20BalanceOf(data.Token, metadata.MsgSender),
				data.Amount.ToBig(),
			),
		)
		env.DelegateCallVoucher(safeERC20TransferAddress, delegateCallPayload)
		return nil

	case "safe_transfer_targeted":
		var data struct {
			Token  common.Address `json:"token" validate:"required"`
			To     common.Address `json:"to" validate:"required"`
			Amount *uint256.Int   `json:"amount" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		if err := validator.Struct(data); err != nil {
			return fmt.Errorf("failed to validate input: %w", err)
		}
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
		abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return err
		}
		safeTransferTargetedPayload, err := abiInterface.Pack(
			"safeTransferTargeted",
			data.Token,
			data.To,
			data.To,
			data.Amount.ToBig(),
		)
		if err != nil {
			return err
		}
		env.SetERC20Balance(
			data.Token,
			metadata.MsgSender,
			new(big.Int).Sub(
				env.ERC20BalanceOf(data.Token, metadata.MsgSender),
				data.Amount.ToBig(),
			),
		)
		env.DelegateCallVoucher(safeERC20TransferAddress, safeTransferTargetedPayload)
		return nil

	case "emergency_erc20_withdraw":
		var data struct {
			Token common.Address `json:"token" validate:"required"`
			To    common.Address `json:"to" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		if err := validator.Struct(data); err != nil {
			return fmt.Errorf("failed to validate input: %w", err)
		}
		abiJSON := `[{
			"type":"function",
			"name":"emergencyERC20Withdraw",
			"inputs":[
				{"type":"address"},
				{"type":"address"}
			]
		}]`
		abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return err
		}
		emergencyERC20WithdrawPayload, err := abiInterface.Pack(
			"emergencyERC20Withdraw",
			data.Token,
			data.To,
		)
		if err != nil {
			return err
		}
		env.DelegateCallVoucher(emergencyWithdrawAddress, emergencyERC20WithdrawPayload)
		return nil

	case "emergency_eth_withdraw":
		var data struct {
			To common.Address `json:"to" validate:"required"`
		}
		if err := json.Unmarshal(input.Data, &data); err != nil {
			return err
		}
		if err := validator.Struct(data); err != nil {
			return fmt.Errorf("failed to validate input: %w", err)
		}
		abiJSON := `[{
			"type":"function",
			"name":"emergencyETHWithdraw",
			"inputs":[
				{"type":"address"}
			]
		}]`
		abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
		if err != nil {
			return err
		}
		emergencyETHWithdrawPayload, err := abiInterface.Pack(
			"emergencyETHWithdraw",
			data.To,
		)
		if err != nil {
			return err
		}
		env.DelegateCallVoucher(emergencyWithdrawAddress, emergencyETHWithdrawPayload)
		return nil

	default:
		env.Report([]byte(fmt.Sprintf("Unknown path: %s", input.Path)))
		return fmt.Errorf("unknown path: %s", input.Path)
	}
}

func (a *Application) Inspect(env rollmelette.EnvInspector, payload []byte) error {
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
	case "contracts":
		contractsJson := fmt.Sprintf(
			`[{"name":"NFT","address":"%s"},{"name":"NFTFactory","address":"%s"},{"name":"EmergencyWithdraw","address":"%s"},{"name":"SafeERC20Transfer","address":"%s"}]`,
			nftAddress,
			nftFactoryAddress,
			emergencyWithdrawAddress,
			safeERC20TransferAddress,
		)
		env.Report([]byte(contractsJson))
		return nil
	default:
		env.Report([]byte(fmt.Sprintf("Unknown path: %s", input.Path)))
		return fmt.Errorf("unknown path: %s", input.Path)
	}
}

func deployNFT(initialOwner common.Address, salt common.Hash, name string, symbol string) ([]byte, error) {
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

func mintNFT(nft common.Address, to common.Address, uri string) ([]byte, error) {
	abiJSON := `[{
		"type":"function",
		"name":"safeMint",
		"inputs":[
			{"type":"address"},
			{"type":"address"},
			{"type":"string"}
		]
	}]`
	abiInterface, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	safeMintPayload, err := abiInterface.Pack(
		"safeMint",
		nft,
		to,
		uri,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI: %w", err)
	}
	return safeMintPayload, nil
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
