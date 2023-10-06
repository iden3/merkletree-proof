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
	"github.com/iden3/merkletree-proof/common"
)

type EthRpcReverseHashCli struct {
	Client             *ethclient.Client
	Contract           *abi.IRHSStorage
	txOpts             *bind.TransactOpts
	RPCResponseTimeout time.Duration
}

func NewEthRpcReverseHashCli(
	contractAddress ethcommon.Address,
	ethClient *ethclient.Client,
	txOpts *bind.TransactOpts,
	timeout time.Duration,
) (*EthRpcReverseHashCli, error) {
	contract, err := abi.NewIRHSStorage(contractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}

	return &EthRpcReverseHashCli{
		Client:             ethClient,
		Contract:           contract,
		txOpts:             txOpts,
		RPCResponseTimeout: timeout,
	}, nil
}

func (cli *EthRpcReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {

	return common.GenerateProof(ctx, cli, treeRoot, key)
}

func (cli *EthRpcReverseHashCli) GetNode(ctx context.Context, hash *merkletree.Hash) (common.Node, error) {
	id := hash.BigInt()

	children, err := cli.Contract.GetNode(nil, id)
	if err != nil {
		if strings.Contains(err.Error(), "Node not found") {
			return common.Node{}, errors.New("node not found")
		}
		return common.Node{}, err
	}

	childrenHashes := make([]*merkletree.Hash, len(children))
	for i, child := range children {
		b, _ := merkletree.NewHashFromBigInt(child)
		childrenHashes[i] = b
	}

	return common.Node{
		Hash:     hash,
		Children: childrenHashes,
	}, nil
}

func (cli *EthRpcReverseHashCli) SaveNodes(ctx context.Context,
	nodes []common.Node) error {

	ctxWT, cancel := context.WithTimeout(ctx, cli.RPCResponseTimeout)
	defer cancel()
	// TODO check if everything is here
	txOpts := &bind.TransactOpts{
		From:      cli.txOpts.From,
		Signer:    cli.txOpts.Signer,
		GasFeeCap: cli.txOpts.GasFeeCap,
		GasTipCap: cli.txOpts.GasTipCap,
		Context:   ctxWT,
		NoSend:    false,
	}

	nodesBigInt := make([][]*big.Int, len(nodes))
	for i, node := range nodes {
		nodesBigInt[i] = make([]*big.Int, len(node.Children))
		for j, child := range node.Children {
			nodesBigInt[i][j] = child.BigInt()
		}
	}

	_, err := cli.Contract.SaveNodes(txOpts, nodesBigInt)
	if err != nil {
		return err
	}

	return nil
}
