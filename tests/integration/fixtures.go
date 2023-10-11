package integration

import (
	"context"
	"encoding/hex"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/merkletree-proof/eth"
)

func NewTestSigner() *Signer {
	pk, _ := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	return &Signer{
		PrivateKey: pk,
		ChainId:    big.NewInt(31337),
	}
}

func NewTestEthRpcReserveHashCli(contractAddress common.Address) (*eth.EthRpcReverseHashCli, error) {
	ethCl, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	timeout := 10 * time.Second
	signer := NewTestSigner()

	addr, _ := signer.Address()

	ctx := context.Background()
	ctxWT, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	tip, err := suggestGasTipCap(ctxWT, ethCl, err)
	if err != nil {
		return nil, err
	}

	ctxWT2, cancel2 := context.WithTimeout(ctx, timeout)
	defer cancel2()
	txOpts := &bind.TransactOpts{
		From:      addr,
		Signer:    signer.SignerFn,
		GasTipCap: tip, // The only option we need to set is gasTipCap as some Ethereum nodes don't support eth_maxPriorityFeePerGas
		GasLimit:  0,   // go-ethereum library will estimate gas limit automatically if it is 0
		Context:   ctxWT2,
		NoSend:    false,
	}

	return eth.NewEthRpcReverseHashCli(contractAddress, ethCl, txOpts, timeout)
}

func suggestGasTipCap(ctx context.Context, ethCl *ethclient.Client, err error) (*big.Int, error) {
	tip, err := ethCl.SuggestGasTipCap(ctx)
	// since hardhat doesn't support 'eth_maxPriorityFeePerGas' rpc call.
	// we should hard code 0 as a mainer tips. More information: https://github.com/NomicFoundation/hardhat/issues/1664#issuecomment-1149006010
	if err != nil && strings.Contains(err.Error(), "eth_maxPriorityFeePerGas not found") {
		tip = big.NewInt(0)
	} else if err != nil {
		return nil, err
	}
	return tip, nil
}
