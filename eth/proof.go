package eth

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	misc "github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/common"
	"github.com/iden3/merkletree-proof/eth/abi"
)

type EthRpcReverseHashCli struct {
	Config           *ClientConfig
	Client           *ethclient.Client
	OnChainTreeStore *abi.OnchainIdentityTreeStore
	Signer           Signer
}

func NewEthRpcReverseHashCli(onChainTreeAddress string, ethereumNodeURL string, signer Signer) (*EthRpcReverseHashCli, error) {
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
		Signer:           signer,
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
	nodes []*big.Int) error {

	// new transaction options
	txOpts := &bind.TransactOpts{
		From:      ethcommon.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		Nonce:     nil,
		Signer:    c.Signer.SignerFn(),
		GasPrice:  nil,
		GasFeeCap: c.Config.MaxGasPrice,
		GasTipCap: c.Config.MinGasPrice,
		//GasLimit:  uint64(c.Config.DefaultGasLimit),
		Context: ctx,
		NoSend:  false,
	}

	addNodeTx, err := c.OnChainTreeStore.AddNode(txOpts, nodes)
	if err != nil {
		return err
	}

	fmt.Println("addNodeTx", addNodeTx.Hash().Hex())

	return nil
}

func (c *EthRpcReverseHashCli) CreateRawTx(ctx context.Context, txParams TransactionParams) (*types.Transaction, error) {
	if txParams.Nonce == nil {
		_ctx, cancel := context.WithTimeout(ctx, c.Config.RPCResponseTimeout)
		defer cancel()
		nonce, err := c.Client.PendingNonceAt(_ctx, txParams.FromAddress)
		if err != nil {
			return nil, errors.New("failed to get nonce")
		}
		txParams.Nonce = &nonce
	}

	_ctx2, cancel2 := context.WithTimeout(ctx, c.Config.RPCResponseTimeout)
	defer cancel2()
	gasLimit, err := c.Client.EstimateGas(_ctx2, ethereum.CallMsg{
		From:  txParams.FromAddress, // the sender of the 'transaction'
		To:    &txParams.ToAddress,
		Gas:   0,             // wei <-> gas exchange ratio
		Value: big.NewInt(0), // amount of wei sent along with the call
		Data:  txParams.Payload,
	})
	if err != nil {
		return nil, errors.New("failed to estimate gas")
	}

	latestBlockHeader, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	if txParams.BaseFee == nil {
		// since ETH and Polygon blockchain already supports London fork.
		// no need set special block.
		baseFee := misc.CalcBaseFee(&params.ChainConfig{LondonBlock: big.NewInt(1)}, latestBlockHeader)

		// add 25% to baseFee. baseFee always small value.
		// since we use dynamic fee transactions we will get not used gas back.
		b := math.Round(float64(baseFee.Int64()) * 1.25)
		baseFee = big.NewInt(int64(b))
		txParams.BaseFee = baseFee
	}

	if txParams.GasTips == nil {
		_ctx3, cancel3 := context.WithTimeout(ctx, c.Config.RPCResponseTimeout)
		defer cancel3()
		gasTip, err := c.Client.SuggestGasTipCap(_ctx3)
		// since hardhad doesn't support 'eth_maxPriorityFeePerGas' rpc call.
		// we should hardcode 0 as a mainer tips. More information: https://github.com/NomicFoundation/hardhat/issues/1664#issuecomment-1149006010
		if err != nil && strings.Contains(err.Error(), "eth_maxPriorityFeePerGas not found") {
			fmt.Println("failed get suggest gas tip: %s. use 0 instead", err)
			gasTip = big.NewInt(0)
		} else if err != nil {
			return nil, errors.New("failed to get suggest gas tip")
		}
		txParams.GasTips = gasTip
	}

	maxGasPricePerFee := big.NewInt(0).Add(txParams.BaseFee, txParams.GasTips)
	baseTx := &types.DynamicFeeTx{
		To:        &txParams.ToAddress,
		Nonce:     *txParams.Nonce,
		Gas:       gasLimit,
		Value:     big.NewInt(0),
		Data:      txParams.Payload,
		GasTipCap: txParams.GasTips,
		GasFeeCap: maxGasPricePerFee,
	}

	tx := types.NewTx(baseTx)

	return tx, nil
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
