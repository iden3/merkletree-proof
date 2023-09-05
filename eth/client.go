package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO keep it as it is?
type TransactionParams struct {
	BaseFee     *big.Int
	GasTips     *big.Int
	Nonce       *uint64
	FromAddress common.Address
	ToAddress   common.Address
	Payload     []byte
}

type ClientConfig struct {
	ReceiptTimeout         time.Duration `json:"receipt_timeout"`
	ConfirmationTimeout    time.Duration `json:"confirmation_timeout"`
	ConfirmationBlockCount int64         `json:"confirmation_block_count"`
	DefaultGasLimit        int           `json:"default_gas_limit"`
	MinGasPrice            *big.Int      `json:"min_gas_price"`
	MaxGasPrice            *big.Int      `json:"max_gas_price"`
	RPCResponseTimeout     time.Duration `json:"rpc_response_time_out"`
	WaitReceiptCycleTime   time.Duration `json:"wait_receipt_cycle_time_out"`
	WaitBlockCycleTime     time.Duration `json:"wait_block_cycle_time_out"`
}

// TODO remove excessive methods
type Client interface {
	BalanceAt(ctx context.Context, addr common.Address) (*big.Int, error)
	Call(fn func(*ethclient.Client) error) error
	CurrentBlock(ctx context.Context) (*big.Int, error)
	ChainID(ctx context.Context) (*big.Int, error)
	GetTransactionReceiptByID(ctx context.Context, txID string) (*types.Receipt, error)
	WaitTransactionReceiptByID(ctx context.Context, txID string) (*types.Receipt, error)
	GetTransactionByID(ctx context.Context, txID string) (*types.Transaction, bool, error)
	CreateRawTx(ctx context.Context, txParams TransactionParams) (*types.Transaction, error)
	SendRawTx(ctx context.Context, tx *types.Transaction) error
	getGasPrice(ctx context.Context) (*big.Int, error)
}

type Signer interface {
	Sign([]byte) ([]byte, error)
	SignerFn() func(common.Address, *types.Transaction) (*types.Transaction, error)
}
