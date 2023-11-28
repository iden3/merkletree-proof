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
	from                 ethcommon.Address
	signer               bind.SignerFn
	rpcTimeout           time.Duration
	needWaitReceipt      bool
	txReceiptTimeout     time.Duration
	waitReceiptCycleTime time.Duration
}

type Option func(cli *ReverseHashCli) error

func WithEthClient(ethClient *ethclient.Client) Option {
	return func(cli *ReverseHashCli) error {
		cli.ethClient = ethClient
		return nil
	}
}

func WithRPCTimeout(timeout time.Duration) Option {
	return func(cli *ReverseHashCli) error {
		cli.rpcTimeout = timeout
		return nil
	}
}

func WithNeedWaitReceipt(needWaitReceipt bool) Option {
	return func(cli *ReverseHashCli) error {
		cli.needWaitReceipt = needWaitReceipt
		return nil
	}
}

func WithTxReceiptTimeout(timeout time.Duration) Option {
	return func(cli *ReverseHashCli) error {
		cli.txReceiptTimeout = timeout
		return nil
	}
}

func WithWaitReceiptCycleTime(cycleTime time.Duration) Option {
	return func(cli *ReverseHashCli) error {
		cli.waitReceiptCycleTime = cycleTime
		return nil
	}
}

func NewReverseHashCli(contractAddress ethcommon.Address,
	from ethcommon.Address, signerFn bind.SignerFn,
	opts ...Option) (*ReverseHashCli, error) {

	ethCl, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		return nil, err
	}

	rhc := &ReverseHashCli{
		ethClient:            ethCl,
		from:                 from,
		signer:               signerFn,
		rpcTimeout:           30 * time.Second,
		needWaitReceipt:      false,
		txReceiptTimeout:     30 * time.Second,
		waitReceiptCycleTime: time.Second,
	}

	for _, o := range opts {
		err := o(rhc)
		if err != nil {
			return nil, err
		}
	}

	contract, err := abi.NewIRHSStorage(contractAddress, rhc.ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}
	rhc.contract = contract

	return rhc, nil
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

	nodesBigInt := make([][]*big.Int, len(nodes))
	for i, node := range nodes {
		nodesBigInt[i] = make([]*big.Int, len(node.Children))
		for j, child := range node.Children {
			nodesBigInt[i][j] = child.BigInt()
		}
	}

	ctxRPC, cancelRPC := cli.ctxWithRPCTimeout(ctx)
	defer cancelRPC()

	txOpts, err := cli.txOptions(ctx, ctxRPC)
	if err != nil {
		return err
	}

	tx, err := cli.contract.SaveNodes(txOpts, nodesBigInt)
	if err != nil {
		return err
	}

	_, err = cli.waitReceipt(ctx, cli.ethClient, tx)
	if err != nil {
		return err
	}
	return nil
}

func (cli *ReverseHashCli) txOptions(ctx, ctxRPC context.Context) (*bind.TransactOpts, error) {
	gasTipCap, err := cli.suggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	txOpts := &bind.TransactOpts{
		From:      cli.from,
		Signer:    cli.signer,
		GasTipCap: gasTipCap,
		GasLimit:  0, // go-ethereum library will estimate gas limit automatically if it is 0
		Context:   ctxRPC,
		NoSend:    false,
	}
	return txOpts, nil
}

func (cli *ReverseHashCli) ctxWithRPCTimeout(ctx context.Context) (context.Context,
	context.CancelFunc) {

	if cli.rpcTimeout > 0 {
		return context.WithTimeout(cli.ctx(ctx), cli.rpcTimeout)
	}

	return ctx, func() {}
}

func (cli *ReverseHashCli) ctxWithTxReceiptTimeout(ctx context.Context) (context.Context,
	context.CancelFunc) {

	if cli.txReceiptTimeout > 0 {
		return context.WithTimeout(cli.ctx(ctx), cli.txReceiptTimeout)
	}

	return ctx, func() {}
}

func (cli *ReverseHashCli) ctx(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return ctx
}

func (cli *ReverseHashCli) waitReceipt(ctx context.Context,
	cl *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {

	if !cli.needWaitReceipt {
		return nil, nil
	}

	ctx, cancel := cli.ctxWithTxReceiptTimeout(ctx)
	defer cancel()

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
			return nil, err
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
