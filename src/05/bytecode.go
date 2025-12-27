package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

//go:embed NFT.json
var NFTJson []byte

type NFTArtifact struct {
	Bytecode string `json:"bytecode"`
}

func getNFTBytecode() ([]byte, error) {
	var artifact NFTArtifact
	if err := json.Unmarshal(NFTJson, &artifact); err != nil {
		return nil, fmt.Errorf("failed to parse embedded NFT.json: %w", err)
	}

	if artifact.Bytecode == "" {
		return nil, fmt.Errorf("bytecode not found in embedded Badge.json")
	}

	bytecode := common.Hex2Bytes(strings.TrimPrefix(artifact.Bytecode, "0x"))

	return bytecode, nil
}
