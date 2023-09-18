package eth

import (
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

type Signer struct {
	PrivateKey []byte
	ChainId    *big.Int
}

func (s *Signer) SignerFn(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(s.ChainId)
	h := signer.Hash(tx)
	sig, err := s.sign(h[:])
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(signer, sig)
}

func (s *Signer) Address() (common.Address, error) {
	pubKey, err := s.getPublicKey()
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(pubKey), nil
}

func (s *Signer) sign(data []byte) ([]byte, error) {
	privKey, err := s.getPrivateKey()
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(data, privKey)
	return sig, err
}

func (s *Signer) getPrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(s.PrivateKey)
}

func (s *Signer) getPublicKey() (ecdsa.PublicKey, error) {
	privKey, err := s.getPrivateKey()
	if err != nil {
		return ecdsa.PublicKey{}, err
	}
	return privKey.PublicKey, nil
}
