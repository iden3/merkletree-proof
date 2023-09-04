package integration

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/common"
	"github.com/iden3/merkletree-proof/eth"
	"github.com/stretchr/testify/assert"
)

func createMockEthRpcReverseHashCli() *eth.EthRpcReverseHashCli {
	cli, err := eth.NewEthRpcReverseHashCli("0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "")
	if err != nil {
		panic(err)
	}
	return cli
}

func TestEthRpcReverseHashCli_GenerateProof(t *testing.T) {
	cli := createMockEthRpcReverseHashCli()
	treeRoot := &merkletree.Hash{}
	key := &merkletree.Hash{}

	proof, err := cli.GenerateProof(context.Background(), treeRoot, key)

	assert.NoError(t, err)
	assert.NotNil(t, proof)
}

func TestEthRpcReverseHashCli_GetNode(t *testing.T) {
	cli := createMockEthRpcReverseHashCli()

	id, _ := big.NewInt(0).SetString("19392314395028218855071922567043158305035792433175725594195224138645494498149", 10)
	node, err := cli.GetNode(context.Background(), id)

	fmt.Println(node)

	assert.NoError(t, err)
	assert.NotNil(t, node)
}

func TestEthRpcReverseHashCli_SaveNodes(t *testing.T) {
	cli := createMockEthRpcReverseHashCli()
	nodes := []common.Node{}

	err := cli.SaveNodes(context.Background(), nodes)

	assert.NoError(t, err)
}
