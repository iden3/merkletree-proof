package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/contracts-abi/rhs-storage/go/abi"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof"
)

type ReverseHashCli struct {
	contract   *abi.IRHSStorage
	txOpts     *bind.TransactOpts
	rpcTimeout time.Duration
}

func NewReverseHashCli(contractAddress ethcommon.Address,
	ethClient *ethclient.Client, txOpts *bind.TransactOpts,
	defaultRPCTimeout time.Duration) (*ReverseHashCli, error) {

	if ethClient == nil {
		return nil, errors.New("ethClient is nil")
	}

	contract, err := abi.NewIRHSStorage(contractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}

	return &ReverseHashCli{
		contract:   contract,
		txOpts:     txOpts,
		rpcTimeout: defaultRPCTimeout,
	}, nil
}

func (cli *ReverseHashCli) ctx(
	ctx context.Context) (context.Context, context.CancelFunc) {

	if ctx == nil {
		ctx = cli.txOpts.Context
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if cli.rpcTimeout > 0 {
		return context.WithTimeout(ctx, cli.rpcTimeout)
	}

	return ctx, func() {}
}

func (cli *ReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {

	return merkletree_proof.GenerateProof(ctx, cli, treeRoot, key)
}

func (cli *ReverseHashCli) GetNode(ctx context.Context,
	hash *merkletree.Hash) (merkletree_proof.Node, error) {

	id := hash.BigInt()

	ctx, cancel := cli.ctx(ctx)
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

	ctx, cancel := cli.ctx(ctx)
	defer cancel()

	txOpts := &bind.TransactOpts{
		From:      cli.txOpts.From,
		Signer:    cli.txOpts.Signer,
		GasFeeCap: cli.txOpts.GasFeeCap,
		GasTipCap: cli.txOpts.GasTipCap,
		Context:   ctx,
		NoSend:    false,
	}

	nodesBigInt := make([][]*big.Int, len(nodes))
	for i, node := range nodes {
		nodesBigInt[i] = make([]*big.Int, len(node.Children))
		for j, child := range node.Children {
			nodesBigInt[i][j] = child.BigInt()
		}
	}

	_, err := cli.contract.SaveNodes(txOpts, nodesBigInt)
	if err != nil {
		return err
	}

	return nil
}
