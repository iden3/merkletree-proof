package common

import (
	"context"
	"encoding/json"
	"fmt"

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
