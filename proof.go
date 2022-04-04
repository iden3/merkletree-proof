package merkletree_proof

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/iden3/go-iden3-crypto/constants"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"
)

const (
	jsonKeyHashKey   = "key"
	jsonKeyHashValue = "value"
	jsonKeyNodeAux   = "aux_node"
	jsonKeyExistence = "existence"
	jsonKeySiblings  = "siblings"
)

func init() {
	var err error
	hashOne, err = NewHashFromBigInt(big.NewInt(1))
	if err != nil {
		panic(err)
	}
}

type Hash [32]byte

var hashZero = Hash{}
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

type Proof struct {
	Existence bool
	Siblings  []Hash
	NodeAux   *LeafNode
}

func (p *Proof) UnmarshalJSON(data []byte) error {
	var obj map[string]json.RawMessage
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	raw, ok := obj[jsonKeyExistence]
	if !ok {
		return fmt.Errorf("missing '%v' key", jsonKeyExistence)
	}
	err = json.Unmarshal(raw, &p.Existence)
	if err != nil {
		return err
	}

	raw, ok = obj[jsonKeySiblings]
	if ok {
		err = json.Unmarshal(raw, &p.Siblings)
		if err != nil {
			return err
		}
	} else {
		p.Siblings = nil
	}

	raw, ok = obj[jsonKeyNodeAux]
	if ok {
		err = json.Unmarshal(raw, &p.NodeAux)
		if err != nil {
			return err
		}
	} else {
		p.NodeAux = nil
	}

	return nil
}

func (p Proof) MarshalJSON() ([]byte, error) {
	siblings := make([]string, len(p.Siblings))
	for i := range p.Siblings {
		siblings[i] = p.Siblings[i].Hex()
	}
	obj := map[string]interface{}{
		jsonKeyExistence: p.Existence,
		jsonKeySiblings:  siblings}
	if p.NodeAux != nil {
		obj[jsonKeyNodeAux] = p.NodeAux
	}
	return json.Marshal(obj)
}

// Root calculates tree root
func (p Proof) Root(key, value Hash) (Hash, error) {
	var midKey *big.Int
	var err error

	if p.Existence {
		midKey, err = leafKey(key.Int(), value.Int())
		if err != nil {
			return Hash{}, err
		}
	} else {
		if p.NodeAux == nil {
			midKey = constants.Zero
		} else {
			if key == p.NodeAux.Key {
				return Hash{}, errors.New(
					"non-existence proof being checked against hIndex equal " +
						"to nodeAux")
			}
			midKey, err = leafKey(p.NodeAux.Key.Int(), p.NodeAux.Value.Int())
			if err != nil {
				return Hash{}, err
			}
		}
	}

	for lvl := len(p.Siblings) - 1; lvl >= 0; lvl-- {
		var left, right *big.Int
		if testBitLittleEndian(key[:], uint(lvl)) {
			left = p.Siblings[lvl].Int()
			right = midKey
		} else {
			left = midKey
			right = p.Siblings[lvl].Int()
		}
		midKey, err = middleNodeKey(left, right)
		if err != nil {
			return Hash{}, err
		}
	}

	return NewHashFromBigInt(midKey)
}

// calculates hash of leaf node
func leafKey(k, v *big.Int) (*big.Int, error) {
	return poseidon.Hash([]*big.Int{k, v, constants.One})
}

// calculates hash of middle node
func middleNodeKey(left, right *big.Int) (*big.Int, error) {
	return poseidon.Hash([]*big.Int{left, right})
}

type NodeType byte

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeMiddle  NodeType = iota
	NodeTypeLeaf    NodeType = iota
	NodeTypeState   NodeType = iota
)

var ErrNodeNotFound = errors.New("node not found")

// GenerateProof generates proof of existence or in-existence of a key in
// a tree identified by a treeRoot.
func GenerateProof(rhsURL string, treeRoot Hash, key Hash) (Proof, error) {
	nextKey := treeRoot
	var p Proof
	for depth := uint(0); depth < uint(len(key)*8); depth++ {
		if nextKey == hashZero {
			return p, nil
		}
		n, err := getNodeFromRHS(rhsURL, nextKey)
		if err != nil {
			return p, err
		}
		switch nt := nodeType(n); nt {
		case NodeTypeLeaf:
			if key == n.Children[0] {
				p.Existence = true
				return p, nil
			}
			// We found a leaf whose entry didn't match hIndex
			p.NodeAux = &LeafNode{Key: n.Children[0], Value: n.Children[1]}
			return p, nil
		case NodeTypeMiddle:
			var siblingKey Hash
			if testBitLittleEndian(key[:], depth) {
				nextKey = n.Children[1]
				siblingKey = n.Children[0]
			} else {
				nextKey = n.Children[0]
				siblingKey = n.Children[1]
			}
			p.Siblings = append(p.Siblings, siblingKey)
		default:
			return p, fmt.Errorf(
				"found unexpected node type in tree (%v): %v",
				nt, n.Hash.Hex())
		}
	}

	return p, errors.New("tree depth is too high")
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

func getNodeFromRHS(rhsURL string, hash Hash) (Node, error) {
	rhsURL = strings.TrimSuffix(rhsURL, "/")
	rhsURL += "/node/" + hash.Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, rhsURL, http.NoBody)
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

	// begin debug
	// x, err := json.Marshal(nodeResp.Node)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf(string(x))
	// end debug

	return nodeResp.Node, nil
}

// testBitLittleEndian tests whether the bit n in bitmap is 1.
func testBitLittleEndian(bitmap []byte, n uint) bool {
	return bitmap[n/8]&(1<<(n%8)) != 0
}
