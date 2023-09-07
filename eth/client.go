package eth

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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

type CliSigner interface {
	Address() (common.Address, error)
	SignerFn(common.Address, *types.Transaction) (*types.Transaction, error)
}
