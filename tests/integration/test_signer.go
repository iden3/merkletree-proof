package integration

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type TestSigner struct {
	PrivateKey []byte
	ChainId    *big.Int
}

func (s *TestSigner) SignerFn(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(s.ChainId)
	h := signer.Hash(tx)
	sig, err := s.sign(h[:])
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(signer, sig)
}

func (s *TestSigner) Address() (common.Address, error) {
	pubKey, err := s.getPublicKey()
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(pubKey), nil
}

func (s *TestSigner) sign(data []byte) ([]byte, error) {
	privKey, err := s.getPrivateKey()
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(data, privKey)
	return sig, err
}

func (s *TestSigner) getPrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(s.PrivateKey)
}

func (s *TestSigner) getPublicKey() (ecdsa.PublicKey, error) {
	privKey, err := s.getPrivateKey()
	if err != nil {
		return ecdsa.PublicKey{}, err
	}
	return privKey.PublicKey, nil
}
