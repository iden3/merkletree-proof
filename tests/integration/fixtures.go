package integration

import (
	"encoding/hex"
	"math/big"
	"time"

	"github.com/iden3/merkletree-proof/eth"
)

func NewTestSigner() *eth.Signer {
	pk, _ := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	return &eth.Signer{
		PrivateKey: pk,
		ChainId:    big.NewInt(31337),
	}
}

func NewTestEthRpcReserveHashCli(contractAddress string, signer eth.CliSigner) (*eth.EthRpcReverseHashCli, error) {
	config := &eth.ClientConfig{
		ReceiptTimeout:         5 * time.Second,
		ConfirmationTimeout:    10 * time.Second,
		ConfirmationBlockCount: 6,
		DefaultGasLimit:        21000,
		MinGasPrice:            big.NewInt(1000000000),
		MaxGasPrice:            big.NewInt(2000000000),
		RPCResponseTimeout:     5 * time.Second,
		WaitReceiptCycleTime:   1 * time.Second,
		WaitBlockCycleTime:     1 * time.Second,
	}
	rpcUrl := "http://127.0.0.1:8545"

	return eth.NewEthRpcReverseHashCli(contractAddress, rpcUrl, signer, config)
}
