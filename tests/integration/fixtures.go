package integration

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/merkletree-proof/eth"
)

func NewTestEthRpcReserveHashCli(contractAddress common.Address) (*eth.ReverseHashCli, error) {
	ethCl, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}
	signer := NewTestSigner()
	fromAddr, _ := signer.Address()

	return eth.NewReverseHashCli(ethCl,
		contractAddress,
		fromAddr,
		signer.SignerFn,
	)
}

func NewTestSigner() *TestSigner {
	pk, _ := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	chainID := int64(31337)

	return &TestSigner{
		PrivateKey: pk,
		ChainId:    big.NewInt(chainID),
	}
}
