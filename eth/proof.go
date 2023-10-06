package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/common"
	"github.com/iden3/merkletree-proof/eth/contracts"
)

type EthRpcReverseHashCli struct {
	Config   *ClientConfig
	Client   *ethclient.Client
	Contract *contracts.IdentityTreeStore
	txOpts   *bind.TransactOpts
}

func NewEthRpcReverseHashCli(
	contractAddress ethcommon.Address, ethClient *ethclient.Client, txOpts *bind.TransactOpts,
) (*EthRpcReverseHashCli, error) {
	contract, err := contracts.NewIdentityTreeStore(contractAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}

	return &EthRpcReverseHashCli{
		Client:   ethClient,
		Contract: contract,
		txOpts:   txOpts,
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

	ctxWT, cancel := context.WithTimeout(ctx, cli.Config.RPCResponseTimeout)
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

// HeaderByNumber get eth block by block number
func (cli *EthRpcReverseHashCli) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	_ctx, cancel := context.WithTimeout(ctx, cli.Config.RPCResponseTimeout)
	defer cancel()
	header, err := cli.Client.HeaderByNumber(_ctx, number)
	if err != nil {
		return nil, err
	}
	return header, nil
}