package eth

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type MockSigner struct {
	PrivateKey []byte
	ChainId    *big.Int
}

func (s *MockSigner) Sign(data []byte) ([]byte, error) {
	privKey, err := crypto.ToECDSA(s.PrivateKey)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(data, privKey)
	return sig, err
}

func (s *MockSigner) SignerFn() func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	return func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		signer := types.LatestSignerForChainID(s.ChainId)
		h := signer.Hash(tx)
		sig, err := s.Sign(h[:])
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, sig)
	}
}

func (s *MockSigner) PublicKey() (ecdsa.PublicKey, error) {
	privKey, err := crypto.ToECDSA(s.PrivateKey)
	if err != nil {
		return ecdsa.PublicKey{}, err
	}
	return privKey.PublicKey, nil
}

func (s *MockSigner) Address() (common.Address, error) {
	pubKey, err := s.PublicKey()
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(pubKey), nil
}
