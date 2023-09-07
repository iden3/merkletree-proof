package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/common"
	"github.com/iden3/merkletree-proof/eth/contracts"
)

type EthRpcReverseHashCli struct {
	Config    *ClientConfig
	Client    *ethclient.Client
	Contract  *contracts.OnchainIdentityTreeStore
	CliSigner CliSigner // TODO Consider better naming
}

func NewEthRpcReverseHashCli(onChainTreeAddress string, ethereumNodeURL string, signer CliSigner) (*EthRpcReverseHashCli, error) {
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

	contract, err := contracts.NewOnchainIdentityTreeStore(addr, cl)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a smart contract: %s", err)
	}

	return &EthRpcReverseHashCli{
		Config:    config,
		Client:    cl,
		Contract:  contract,
		CliSigner: signer,
	}, nil
}

func (c *EthRpcReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {
	return &merkletree.Proof{}, nil
}

func (c *EthRpcReverseHashCli) GetNode(ctx context.Context, id *big.Int) (common.Node, error) {
	children, err := c.Contract.GetNode(nil, id)
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
		Children: childrenHashes,
	}, nil
}

func (c *EthRpcReverseHashCli) SaveNodes(ctx context.Context,
	nodes []*big.Int) error {

	addr, err := c.CliSigner.Address()
	if err != nil {
		return err
	}

	// TODO consider if evaluate gas price and hardcap limit is needed
	txOpts := &bind.TransactOpts{
		From:      addr,
		Signer:    c.CliSigner.SignerFn,
		GasFeeCap: c.Config.MaxGasPrice,
		GasTipCap: c.Config.MinGasPrice,
		Context:   ctx,
		NoSend:    false,
	}

	addNodeTx, err := c.Contract.AddNode(txOpts, nodes)
	if err != nil {
		return err
	}

	fmt.Println("addNodeTx", addNodeTx.Hash().Hex())

	return nil
}

// HeaderByNumber get eth block by block number
func (c *EthRpcReverseHashCli) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	_ctx, cancel := context.WithTimeout(ctx, c.Config.RPCResponseTimeout)
	defer cancel()
	header, err := c.Client.HeaderByNumber(_ctx, number)
	if err != nil {
		return nil, err
	}
	return header, nil
}
