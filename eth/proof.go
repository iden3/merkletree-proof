package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/iden3/contracts-abi/rhs-storage/go/abi"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof"
)

type ReverseHashCli struct {
	contract             *abi.IRHSStorage
	ethClient            *ethclient.Client
	txOpts               *bind.TransactOpts
	rpcTimeout           time.Duration
	txReceiptTimeout     time.Duration
	waitReceiptCycleTime time.Duration
}

func NewReverseHashCli(contractAddress ethcommon.Address,
	ethClient *ethclient.Client, txOpts *bind.TransactOpts,
	rpcTimeout time.Duration, txReceiptTimeout time.Duration,
	waitReceiptCycleTime time.Duration) (*ReverseHashCli, error) {

	if ethClient == nil {
		return nil, errors.New("ethClient is nil")
	}

	contract, err := abi.NewIRHSStorage(contractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}

	return &ReverseHashCli{
		contract:             contract,
		ethClient:            ethClient,
		txOpts:               txOpts,
		rpcTimeout:           rpcTimeout,
		txReceiptTimeout:     txReceiptTimeout,
		waitReceiptCycleTime: waitReceiptCycleTime,
	}, nil
}

func (cli *ReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {

	return merkletree_proof.GenerateProof(ctx, cli, treeRoot, key)
}

func (cli *ReverseHashCli) GetNode(ctx context.Context,
	hash *merkletree.Hash) (merkletree_proof.Node, error) {

	id := hash.BigInt()

	ctx, cancel := cli.ctxWithRPCTimeout(ctx)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}
	children, err := cli.contract.GetNode(opts, id)
	if err != nil {
		if strings.Contains(err.Error(), "Node not found") {
			return merkletree_proof.Node{}, merkletree_proof.ErrNodeNotFound
		}
		return merkletree_proof.Node{}, err
	}

	n := merkletree_proof.Node{
		Hash:     hash,
		Children: make([]*merkletree.Hash, len(children)),
	}
	for i, child := range children {
		n.Children[i], err = merkletree.NewHashFromBigInt(child)
		if err != nil {
			return merkletree_proof.Node{}, err
		}
	}

	return n, nil
}

func (cli *ReverseHashCli) SaveNodes(ctx context.Context,
	nodes []merkletree_proof.Node) error {

	gasTipCap, err := cli.suggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	ctxRPC, cancelRPC := cli.ctxWithRPCTimeout(ctx)
	defer cancelRPC()

	txOpts := &bind.TransactOpts{
		From:      cli.txOpts.From,
		Signer:    cli.txOpts.Signer,
		GasFeeCap: cli.txOpts.GasFeeCap,
		GasTipCap: gasTipCap,
		Context:   ctxRPC,
		NoSend:    false,
	}

	nodesBigInt := make([][]*big.Int, len(nodes))
	for i, node := range nodes {
		nodesBigInt[i] = make([]*big.Int, len(node.Children))
		for j, child := range node.Children {
			nodesBigInt[i][j] = child.BigInt()
		}
	}

	tx, err := cli.contract.SaveNodes(txOpts, nodesBigInt)
	if err != nil {
		return err
	}

	ctxRpt, cancelRpt := cli.ctxWithTxReceiptTimeout(ctx)
	defer cancelRpt()

	_, err = cli.waitMined(ctxRpt, cli.ethClient, tx)
	if err != nil {
		return err
	}
	return nil
}

func (cli *ReverseHashCli) ctxWithRPCTimeout(
	ctx context.Context) (context.Context, context.CancelFunc) {
	ctx = cli.ctx(ctx)

	if cli.rpcTimeout > 0 {
		return context.WithTimeout(ctx, cli.rpcTimeout)
	}

	return ctx, func() {}
}

func (cli *ReverseHashCli) ctxWithTxReceiptTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx = cli.ctx(ctx)

	if cli.txReceiptTimeout > 0 {
		return context.WithTimeout(ctx, cli.txReceiptTimeout)
	}

	return ctx, func() {}
}

func (cli *ReverseHashCli) ctx(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = cli.txOpts.Context
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

func (cli *ReverseHashCli) waitMined(ctx context.Context, cl *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	queryTicker := time.NewTicker(cli.waitReceiptCycleTime)
	defer queryTicker.Stop()

	logger := log.New("hash", tx.Hash())
	for {
		receipt, err := cl.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			return receipt, nil
		}

		if errors.Is(err, ethereum.NotFound) {
			logger.Trace("Transaction not yet mined")
		} else {
			logger.Trace("Receipt retrieval failed", "err", err)
		}

		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func (cli *ReverseHashCli) suggestGasTipCap(ctx context.Context) (*big.Int, error) {
	ctxRPC, cancel := cli.ctxWithRPCTimeout(ctx)
	defer cancel()

	tip, err := cli.ethClient.SuggestGasTipCap(ctxRPC)
	// since hardhat doesn't support 'eth_maxPriorityFeePerGas' rpc call.
	// we should hard code 0 as a mainer tips. More information: https://github.com/NomicFoundation/hardhat/issues/1664#issuecomment-1149006010
	if err != nil && strings.Contains(err.Error(), "eth_maxPriorityFeePerGas not found") {
		log.Trace("failed get suggest gas tip. Use 0 instead", "err", err)
		tip = big.NewInt(0)
	} else if err != nil {
		return nil, fmt.Errorf("failed get suggest gas tip: %w", err)
	}

	return tip, nil
}
