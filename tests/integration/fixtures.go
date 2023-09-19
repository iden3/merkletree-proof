package integration

import (
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/merkletree-proof/eth"
)

func NewTestSigner() *eth.Signer {
	pk, _ := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	return &eth.Signer{
		PrivateKey: pk,
		ChainId:    big.NewInt(31337),
	}
}

func NewTestEthRpcReserveHashCli(contractAddress common.Address, signer eth.CliSigner) (*eth.EthRpcReverseHashCli, error) {
	config := &eth.ClientConfig{
		MinGasPrice:        big.NewInt(1000000000),
		MaxGasPrice:        big.NewInt(2000000000),
		RPCResponseTimeout: 5 * time.Second,
	}

	cl, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	return eth.NewEthRpcReverseHashCli(contractAddress, cl, signer, config)
}
