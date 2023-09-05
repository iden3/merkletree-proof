package integration

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/eth"
	"github.com/stretchr/testify/assert"
)

func createMockEthRpcReverseHashCli() *eth.EthRpcReverseHashCli {
	pk, err := hex.DecodeString("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		panic(err)
	}

	signer := &eth.MockSigner{
		PrivateKey: pk,
		ChainId:    big.NewInt(31337),
	}

	cli, err := eth.NewEthRpcReverseHashCli("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "", signer)
	if err != nil {
		panic(err)
	}
	return cli
}

func TestEthRpcReverseHashCli_SaveNodes(t *testing.T) {
	cli := createMockEthRpcReverseHashCli()
	// TODO rewrite to nodes := []common.Node{}
	nodes := make([]*big.Int, 3)
	nodes[0] = big.NewInt(2)
	nodes[1] = big.NewInt(3)
	nodes[2] = big.NewInt(4)

	err := cli.SaveNodes(context.Background(), nodes)

	assert.NoError(t, err)
}

func TestEthRpcReverseHashCli_GetNode(t *testing.T) {
	cli := createMockEthRpcReverseHashCli()

	id, _ := big.NewInt(0).SetString("19392314395028218855071922567043158305035792433175725594195224138645494498149", 10)
	node, err := cli.GetNode(context.Background(), id)

	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, len(node.Children), 3)
}

func TestEthRpcReverseHashCli_GenerateProof(t *testing.T) {
	cli := createMockEthRpcReverseHashCli()
	treeRoot := &merkletree.Hash{}
	key := &merkletree.Hash{}

	proof, err := cli.GenerateProof(context.Background(), treeRoot, key)

	assert.NoError(t, err)
	assert.NotNil(t, proof)
}
