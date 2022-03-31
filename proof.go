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

	"github.com/iden3/go-iden3-crypto/utils"
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

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

type Node struct {
	Hash     Hash
	Children []Hash
}

type NodeAux struct {
	Key   Hash
	Value Hash
}

func (n NodeAux) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		jsonKeyHashKey:   n.Key.Hex(),
		jsonKeyHashValue: n.Value.Hex(),
	})
}

type Proof struct {
	Existence bool     `json:"existence"`
	Siblings  []Hash   `json:"siblings"`
	NodeAux   *NodeAux `json:"aux_node"`
}

func (p *Proof) UnmarshalJSON(data []byte) error {
	p.Siblings = nil
	p.NodeAux = nil

	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}

	exI, ok := obj["existence"]
	if !ok {
		return errors.New("existence key not found")
	}
	p.Existence, ok = exI.(bool)
	if !ok {
		return errors.New("incorrect type of existence key")
	}

	sibI, ok := obj["siblings"]
	if !ok || sibI == nil {
		p.Siblings = nil
	} else {
		sibL, ok := sibI.([]interface{})
		if !ok {
			return fmt.Errorf("incorrect type of siblings key: %T", sibI)
		}
		p.Siblings = make([]Hash, len(sibL))
		for i, s := range sibL {
			sS, ok := s.(string)
			if !ok {
				return fmt.Errorf("sibling #%v is not string", i)
			}
			p.Siblings[i], err = unmarshalHex(sS)
			if err != nil {
				return fmt.Errorf("errors unmarshal sibling #%v: %v", i, err)
			}
		}
	}

	anI, ok := obj["aux_node"]
	if !ok || anI == nil {
		p.NodeAux = nil
		return nil
	}

	anI2, ok := anI.(map[string]interface{})
	if !ok {
		return errors.New("aux_node has incorrect format")
	}

	p.NodeAux = new(NodeAux)

	keyI, ok := anI2["key"]
	if !ok {
		return errors.New("aux_node has not key")
	}

	keyS, ok := keyI.(string)
	if !ok {
		return errors.New("aux_node key is not a string")
	}

	hashBytes, err := hex.DecodeString(keyS)
	if err != nil {
		return err
	}
	if len(hashBytes) != len(p.NodeAux.Key) {
		return errors.New("incorrect aux_node key length")
	}

	copy(p.NodeAux.Key[:], hashBytes)

	valueI, ok := anI2["value"]
	if !ok {
		return errors.New("aux_node has not value")
	}

	valueS, ok := valueI.(string)
	if !ok {
		return errors.New("aux_node value is not a string")
	}

	hashBytes, err = hex.DecodeString(valueS)
	if err != nil {
		return err
	}
	if len(hashBytes) != len(p.NodeAux.Value) {
		return errors.New("incorrect aux_node value length")
	}

	copy(p.NodeAux.Value[:], hashBytes)

	return nil
}

func unmarshalHex(in string) (Hash, error) {
	var h Hash
	data, err := hex.DecodeString(in)
	if err != nil {
		return h, err
	}
	if len(data) != len(h) {
		return h, errors.New("incorrect length")
	}
	copy(h[:], data)
	return h, nil
}

func (p Proof) MarshalJSON() ([]byte, error) {
	siblings := make([]string, len(p.Siblings))
	for i := range p.Siblings {
		siblings[i] = p.Siblings[i].Hex()
	}
	obj := map[string]interface{}{
		"existence": p.Existence,
		"siblings":  siblings}
	if p.NodeAux != nil {
		obj["aux_node"] = p.NodeAux
	}
	return json.Marshal(obj)
}

type NodeType byte

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeMiddle  NodeType = iota
	NodeTypeLeaf    NodeType = iota
	NodeTypeState   NodeType = iota
)

var ErrNodeNotFound = errors.New("node not found")

func generateProof(rhsURL string, treeRoot Hash, key Hash) (Proof, error) {
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
			p.NodeAux = &NodeAux{Key: n.Children[0], Value: n.Children[1]}
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
