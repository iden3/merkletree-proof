package merkletree_proof

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/iden3/go-merkletree-sql/v2"
)

func init() {
	hashOneP, err := merkletree.NewHashFromBigInt(big.NewInt(1))
	if err != nil {
		panic(err)
	}
	copy(hashOne[:], hashOneP[:])
}

var hashOne merkletree.Hash

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

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeMiddle  NodeType = iota
	NodeTypeLeaf    NodeType = iota
	NodeTypeState   NodeType = iota
)

var ErrNodeNotFound = errors.New("node not found")

type HTTPReverseHashCli struct {
	URL         string
	HTTPTimeout time.Duration
}

// GenerateProof generates proof of existence or in-existence of a key in
// a tree identified by a treeRoot.
func (cli *HTTPReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {

	if cli.URL == "" {
		return nil, errors.New("HTTP reverse hash service url is not specified")
	}

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
		switch nt := nodeType(n); nt {
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

func (cli *HTTPReverseHashCli) nodeURL(node *merkletree.Hash) string {
	nodeURL := cli.baseURL() + "/node"
	if node == nil {
		return nodeURL
	}
	return nodeURL + "/" + node.Hex()
}

func (cli *HTTPReverseHashCli) baseURL() string {
	return strings.TrimSuffix(cli.URL, "/")
}

func (cli *HTTPReverseHashCli) getHttpTimeout() time.Duration {
	if cli.HTTPTimeout == 0 {
		return 10 * time.Second
	}
	return cli.HTTPTimeout
}

func (cli *HTTPReverseHashCli) GetNode(ctx context.Context,
	hash *merkletree.Hash) (Node, error) {

	if hash == nil {
		return Node{}, errors.New("hash is nil")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cli.getHttpTimeout())
		defer cancel()
	}

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, cli.nodeURL(hash), http.NoBody)
	if err != nil {
		return Node{}, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return Node{}, err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode == http.StatusNotFound {
		var resp map[string]interface{}
		dec := json.NewDecoder(httpResp.Body)
		err := dec.Decode(&resp)
		if err != nil {
			return Node{}, err
		}
		if resp["status"] == "not found" {
			return Node{}, ErrNodeNotFound
		} else {
			return Node{}, errors.New("unexpected response")
		}
	} else if httpResp.StatusCode != http.StatusOK {
		return Node{}, fmt.Errorf("unexpected response: %v",
			httpResp.StatusCode)
	}

	var nodeResp nodeResponse
	dec := json.NewDecoder(httpResp.Body)
	err = dec.Decode(&nodeResp)
	if err != nil {
		return Node{}, err
	}

	return nodeResp.Node, nil
}

func (cli *HTTPReverseHashCli) SaveNodes(ctx context.Context,
	nodes []Node) error {

	reqBytes, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	// if no timeout set on context, set it here
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cli.getHttpTimeout())
		defer cancel()
	}

	bodyReader := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		cli.nodeURL(nil), bodyReader)
	if err != nil {
		return err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	dec := json.NewDecoder(httpResp.Body)
	var respM map[string]interface{}
	err = dec.Decode(&respM)
	if err != nil {
		return fmt.Errorf("unable to decode RHS response: %w", err)
	}

	if respM["status"] != "OK" {
		return fmt.Errorf("unexpected RHS response status: %s", respM["status"])
	}

	return nil
}

func nodeType(node Node) NodeType {
	if len(node.Children) == 2 {
		return NodeTypeMiddle
	}

	if len(node.Children) == 3 && *node.Children[2] == hashOne {
		return NodeTypeLeaf
	}

	if len(node.Children) == 3 {
		return NodeTypeState
	}

	return NodeTypeUnknown
}

type nodeResponse struct {
	Node   Node   `json:"node"`
	Status string `json:"status"`
}

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
