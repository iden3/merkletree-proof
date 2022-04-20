package merkletree_proof

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-merkletree-sql"
)

const (
	jsonKeyHashKey   = "key"
	jsonKeyHashValue = "value"
)

func init() {
	var err error
	hashOne, err = NewHashFromBigInt(big.NewInt(1))
	if err != nil {
		panic(err)
	}
}

type Hash [32]byte

var hashOne Hash

func NewHashFromBigInt(i *big.Int) (Hash, error) {
	var h Hash
	if !utils.CheckBigIntInField(i) {
		return h, errors.New("big int out of field")
	}
	copy(h[:], utils.SwapEndianness(i.Bytes()))
	return h, nil
}

func NewHashFromHex(in string) (Hash, error) {
	var h Hash

	hashBytes, err := hex.DecodeString(in)
	if err != nil {
		return h, err
	}
	if len(hashBytes) != len(h) {
		return h, errors.New("invalid hash length")
	}

	copy(h[:], hashBytes)
	return h, nil
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) Int() *big.Int {
	return new(big.Int).SetBytes(utils.SwapEndianness(h[:]))
}

func (h *Hash) UnmarshalText(data []byte) error {
	h2, err := NewHashFromHex(string(data))
	if err != nil {
		return err
	}

	if !utils.CheckBigIntInField(h2.Int()) {
		return errors.New("big int out of field")
	}

	copy(h[:], h2[:])
	return nil
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

type Node struct {
	Hash     Hash   `json:"hash"`
	Children []Hash `json:"children"`
}

type LeafNode struct {
	Key   Hash `json:"key"`
	Value Hash `json:"value"`
}

func (n LeafNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		jsonKeyHashKey:   n.Key.Hex(),
		jsonKeyHashValue: n.Value.Hex(),
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
func (cli *HTTPReverseHashCli) GenerateProof(treeRoot *merkletree.Hash,
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
		n, err := cli.GetNode(nextKey)
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
				Key:   hashToMTHash(n.Children[0]),
				Value: hashToMTHash(n.Children[1]),
			}
			return mkProof()
		case NodeTypeMiddle:
			if merkletree.TestBit(key[:], depth) {
				nextKey = hashToMTHash(n.Children[1])
				siblings = append(siblings, hashToMTHash(n.Children[0]))
			} else {
				nextKey = hashToMTHash(n.Children[0])
				siblings = append(siblings, hashToMTHash(n.Children[1]))
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
	nodeURL := cli.baseURL() + "/node/"
	if node == nil {
		return nodeURL
	}
	return nodeURL + node.Hex()
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

func (cli *HTTPReverseHashCli) GetNode(hash *merkletree.Hash) (Node, error) {
	if hash == nil {
		return Node{}, errors.New("hash is nil")
	}

	ctx, cancel :=
		context.WithTimeout(context.Background(), cli.getHttpTimeout())
	defer cancel()

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, cli.nodeURL(hash), http.NoBody)
	if err != nil {
		return Node{}, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return Node{}, err
	}

	defer httpResp.Body.Close()
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

func nodeType(node Node) NodeType {
	if len(node.Children) == 2 {
		return NodeTypeMiddle
	}

	if len(node.Children) == 3 && node.Children[2] == hashOne {
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

func hashToMTHash(hash Hash) *merkletree.Hash {
	var h merkletree.Hash
	copy(h[:], hash[:])
	return &h
}
