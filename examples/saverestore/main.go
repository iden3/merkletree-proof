package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"github.com/iden3/merkletree-proof/common"
	proof "github.com/iden3/merkletree-proof/http"
)

var hashOne *merkletree.Hash

func init() {
	var err error
	hashOne, err = merkletree.NewHashFromBigInt(big.NewInt(1))
	if err != nil {
		panic(err)
	}
}

func main() {
	rhsURL := "http://localhost:8003"
	cli := &proof.HTTPReverseHashCli{URL: rhsURL}

	// generate tree with 10 random leaves.
	tree := buildMT(10)

	// save tree to reverse haash service
	treeNodes := nodesFromTree(tree)
	err := cli.SaveNodes(context.Background(), treeNodes)
	if err != nil {
		panic(err)
	}

	// restore tree from reverse hash service
	rhsTree := restoreTree(cli, tree.Root())

	rhsLeavesDump, err := rhsTree.DumpLeafs(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	origLeavesDump, err := tree.DumpLeafs(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(origLeavesDump, rhsLeavesDump) {
		panic("trees expected to be equal, but not")
	}

	fmt.Println("OK: trees are equal")
}

func restoreTree(cli *proof.HTTPReverseHashCli,
	root *merkletree.Hash) *merkletree.MerkleTree {
	mt := newEmptyTree()
	walkRHSLeafs(cli, root, func(key, value *merkletree.Hash) {
		err := mt.Add(context.Background(), key.BigInt(), value.BigInt())
		if err != nil {
			panic(err)
		}
	})
	return mt
}

func walkRHSLeafs(cli *proof.HTTPReverseHashCli, root *merkletree.Hash,
	fn func(key, value *merkletree.Hash)) {

	node, err := cli.GetNode(context.Background(), root)
	if err != nil {
		panic(err)
	}
	if len(node.Children) == 2 {
		if !node.Children[0].Equals(&merkletree.HashZero) {
			walkRHSLeafs(cli, node.Children[0], fn)
		}
		if !node.Children[1].Equals(&merkletree.HashZero) {
			walkRHSLeafs(cli, node.Children[1], fn)
		}
	} else if len(node.Children) == 3 {
		if !node.Children[2].Equals(hashOne) {
			panic("3rd child of leaf expected to be equal to 1")
		}
		fn(node.Children[0], node.Children[1])
	}
}

func nodesFromTree(tree *merkletree.MerkleTree) []common.Node {
	ctx := context.Background()

	var nodes []common.Node
	err := tree.Walk(ctx, nil, func(node *merkletree.Node) {
		nodeKey, err := node.Key()
		if err != nil {
			panic(err)
		}
		switch node.Type {
		case merkletree.NodeTypeMiddle:
			nodes = append(nodes, common.Node{
				Hash:     nodeKey,
				Children: []*merkletree.Hash{node.ChildL, node.ChildR}})
		case merkletree.NodeTypeLeaf:
			nodes = append(nodes, common.Node{
				Hash: nodeKey,
				Children: []*merkletree.Hash{node.Entry[0], node.Entry[1],
					hashOne},
			})
		case merkletree.NodeTypeEmpty:
			// do not save zero nodes
		default:
			panic(fmt.Sprintf("unexpected node type: %v", node.Type))
		}

	})
	if err != nil {
		panic(err)
	}

	return nodes
}

// n is a number of random leaves in tree
func buildMT(n int) *merkletree.MerkleTree {
	mt := newEmptyTree()
	for _, e := range genRandomEntries(n) {
		err := mt.Add(context.Background(), e.key.BigInt(), e.value.BigInt())
		if err != nil {
			panic(err)
		}
	}
	return mt
}

type entry struct {
	key   *merkletree.Hash
	value *merkletree.Hash
}

// generate random tree entries
func genRandomEntries(n int) []entry {
	result := make([]entry, n)
	for i := 0; i < n; i++ {
		result[i].key = genRandomHash()
		result[i].value = genRandomHash()
	}
	return result
}

// return random *merkletree.Hash
func genRandomHash() *merkletree.Hash {
	var rndValBytes [16]byte
	_, err := rand.Read(rndValBytes[:])
	if err != nil {
		panic(err)
	}
	h, err := merkletree.HashElems(new(big.Int).SetBytes(rndValBytes[:]))
	if err != nil {
		panic(err)
	}
	return h
}

func newEmptyTree() *merkletree.MerkleTree {
	mtStorage := memory.NewMemoryStorage()
	ctx := context.Background()
	const mtDepth = 40
	mt, err := merkletree.NewMerkleTree(ctx, mtStorage, mtDepth)
	if err != nil {
		panic(err)
	}
	return mt
}
