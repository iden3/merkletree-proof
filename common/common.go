package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/iden3/go-merkletree-sql/v2"
)

type ReverseHashCli interface {
	GenerateProof(ctx context.Context,
		treeRoot *merkletree.Hash,
		key *merkletree.Hash) (*merkletree.Proof, error)
	GetNode(ctx context.Context,
		hash *merkletree.Hash) (Node, error)
	SaveNodes(ctx context.Context,
		nodes []Node) error
}

type Node struct {
	Hash     *merkletree.Hash
	Children []*merkletree.Hash
}

type jsonNode struct {
	Hash     string   `json:"hash"`
	Children []string `json:"children"`
}

func (n *Node) UnmarshalJSON(in []byte) error {
	var jsonN jsonNode
	err := json.Unmarshal(in, &jsonN)
	if err != nil {
		return err
	}
	n.Hash, err = merkletree.NewHashFromHex(jsonN.Hash)
	if err != nil {
		return err
	}
	n.Children, err = hexesToHashes(jsonN.Children)
	return err
}

func (n Node) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonNode{
		Hash:     n.Hash.Hex(),
		Children: hashesToHexes(n.Children),
	})
}

type NodeType byte

func hashesToHexes(hashes []*merkletree.Hash) []string {
	if hashes == nil {
		return nil
	}
	hexes := make([]string, len(hashes))
	for i, h := range hashes {
		hexes[i] = h.Hex()
	}
	return hexes
}

func hexesToHashes(hexes []string) ([]*merkletree.Hash, error) {
	if hexes == nil {
		return nil, nil
	}
	hashes := make([]*merkletree.Hash, len(hexes))
	var err error
	for i, h := range hexes {
		hashes[i], err = merkletree.NewHashFromHex(h)
		if err != nil {
			return nil, fmt.Errorf("can't parse hex #%v: %w", i, err)
		}
	}
	return hashes, nil
}

func GenerateProof(ctx context.Context, cli ReverseHashCli,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {

	var exists bool
	var siblings []*merkletree.Hash
	var nodeAux *merkletree.NodeAux

	mkProof := func() (*merkletree.Proof, error) {
		return merkletree.NewProofFromData(exists, siblings, nodeAux)
	}

	nextKey := treeRoot
	for depth := uint(0); depth < uint(len(key)*8); depth++ {
		if *nextKey == merkletree.HashZero {
			return mkProof()
		}
		n, err := cli.GetNode(ctx, nextKey)
		if err != nil {
			return nil, err
		}
		switch nt := GetNodeType(n); nt {
		case NodeTypeLeaf:
			if bytes.Equal(key[:], n.Children[0][:]) {
				exists = true
				return mkProof()
			}
			// We found a leaf whose entry didn't match hIndex
			nodeAux = &merkletree.NodeAux{
				Key:   n.Children[0],
				Value: n.Children[1],
			}
			return mkProof()
		case NodeTypeMiddle:
			if merkletree.TestBit(key[:], depth) {
				nextKey = n.Children[1]
				siblings = append(siblings, n.Children[0])
			} else {
				nextKey = n.Children[0]
				siblings = append(siblings, n.Children[1])
			}
		default:
			return nil, fmt.Errorf(
				"found unexpected node type in tree (%v): %v",
				nt, n.Hash.Hex())
		}
	}

	return nil, errors.New("tree depth is too high")
}

var hashOne, _ = merkletree.NewHashFromBigInt(big.NewInt(1))

func GetNodeType(node Node) NodeType {
	if len(node.Children) == 2 {
		return NodeTypeMiddle
	}

	if len(node.Children) == 3 && *node.Children[2] == *hashOne {
		return NodeTypeLeaf
	}

	if len(node.Children) == 3 {
		return NodeTypeState
	}

	return NodeTypeUnknown
}

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeMiddle  NodeType = iota
	NodeTypeLeaf    NodeType = iota
	NodeTypeState   NodeType = iota
)
