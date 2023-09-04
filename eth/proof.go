package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/common"
	"github.com/iden3/merkletree-proof/eth/abi"
)

type EthRpcReverseHashCli struct {
	Config           *ClientConfig
	Client           *ethclient.Client
	OnChainTreeStore *abi.OnchainIdentityTreeStore
}

func NewEthRpcReverseHashCli(onChainTreeAddress string, ethereumNodeURL string) (*EthRpcReverseHashCli, error) {
	config := &ClientConfig{
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

	url := ethereumNodeURL
	if ethereumNodeURL == "" {
		url = "http://127.0.0.1:8545"
	}

	cl, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %s", err)
	}

	addr := ethcommon.HexToAddress(onChainTreeAddress)

	// connect to the contract
	contract, err := abi.NewOnchainIdentityTreeStore(addr, cl)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}

	return &EthRpcReverseHashCli{
		Config:           config,
		Client:           cl,
		OnChainTreeStore: contract,
	}, nil
}

func (c *EthRpcReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {
	return &merkletree.Proof{}, nil
}

func (c *EthRpcReverseHashCli) GetNode(ctx context.Context, id *big.Int) (common.Node, error) {

	children, err := c.OnChainTreeStore.GetNode(nil, id)
	if err != nil {
		return common.Node{}, err
	}

	childrenHashes := make([]*merkletree.Hash, len(children))
	for i, child := range children {
		b, _ := merkletree.NewHashFromBigInt(child)
		childrenHashes[i] = b
	}

	hash, _ := merkletree.NewHashFromBigInt(id)
	return common.Node{
		Hash:     hash,
		Children: nil,
	}, nil
}

func (c *EthRpcReverseHashCli) SaveNodes(ctx context.Context,
	nodes []common.Node) error {

	return nil
}
