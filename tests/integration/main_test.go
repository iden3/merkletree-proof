package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
	proof "github.com/iden3/merkletree-proof"
	"github.com/stretchr/testify/require"
)

func TestProof(t *testing.T) {
	rhsUrl, ok := os.LookupEnv("RHS_URL")
	if !ok || rhsUrl == "" {
		t.Fatal("RHS_URL not set")
	}

	revNonces := []uint64{
		5577006791947779410,  // 19817761...  0 1 0 0 1 0 1 0
		8674665223082153551,  // 68456430...  1 1 1 1 0 0 1 0
		8674665223082147919,  // a node is very close to 8674665223082153551 — to generate zero siblings
		15352856648520921629, // 86798249...  1 0 1 1 1 0 0 0
		13260572831089785859, // 13668806...  1 1 0 0 0 0 0 0
		3916589616287113937,  // 50401982...  1 0 0 0 1 0 1 1
		6334824724549167320,  // 38589333...  0 0 0 1 1 0 1 1
		9828766684487745566,  // 55091915...  0 1 1 1 1 0 0 0
		10667007354186551956, // 10419680...  0 0 1 0 1 0 0 1
		894385949183117216,   // 13133085...  0 0 0 0 0 1 0 1
		11998794077335055257, // 14875578...  1 0 0 1 1 0 0 1
	}
	bigMerkleTree := buildTree(t, revNonces)
	saveTreeToRHS(t, rhsUrl, bigMerkleTree)

	oneNodeMerkleTree := buildTree(t, []uint64{5577006791947779410})
	saveTreeToRHS(t, rhsUrl, oneNodeMerkleTree)

	t.Run("Test save state", func(t *testing.T) {
		state := saveIdenStateToRHS(t, rhsUrl, bigMerkleTree)

		revTreeRoot, err := getRevTreeRoot(rhsUrl, state)
		require.NoError(t, err)

		revTreeRootExpected := hashFromInt(bigMerkleTree.Root().BigInt())
		require.Equal(t, revTreeRootExpected, revTreeRoot)
	})

	testCases := []struct {
		title       string
		revNonce    uint64
		revTreeRoot proof.Hash
		wantProof   proof.Proof
		wantErr     string
	}{
		{
			title:       "regular node",
			revNonce:    10667007354186551956,
			revTreeRoot: hashFromInt(bigMerkleTree.Root().BigInt()),
			wantProof: proof.Proof{
				Existence: true,
				Siblings: []proof.Hash{
					hashFromHex("74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"),
					hashFromHex("ff11b8bf1d13e28e86e249d2acdba0bd9c0fe4a5f56ad4236b09185bde81c316"),
					hashFromHex("db5eb80f6b60b4e23714d4d00f178ba62fbdb4f0294675f51ac99aa24e600827"),
				},
				NodeAux: nil,
			},
		},
		{
			title:       "a node with zero siblings",
			revNonce:    8674665223082147919,
			revTreeRoot: hashFromInt(bigMerkleTree.Root().BigInt()),
			wantProof: proof.Proof{
				Existence: true,
				Siblings: []proof.Hash{
					hashFromHex("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					hashFromHex("28e5cdd29d9ad96cc214c654ca8e2f4fa5576bc132e172519804a58ee4bb4d18"),
					hashFromHex("658c7a65594ebb0815e1cc20f54284ccdb51bb1625f103c116ce58444145381e"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("e809a4ed2cf98922910e456f1e56862bb958777f5ff0ea6799360113257f220f"),
				},
				NodeAux: nil,
			},
		},
		{
			title: "un-existence with aux node",
			//nolint:gocritic
			revNonce:    5, // revNonceKey[0] = 0b00000101
			revTreeRoot: hashFromInt(bigMerkleTree.Root().BigInt()),
			wantProof: proof.Proof{
				Existence: false,
				Siblings: []proof.Hash{
					hashFromHex("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					hashFromHex("c9719432e3d8bf360d0f2de456c5321c51295895c9330b0588552580765cd929"),
					hashFromHex("c0e8bf477403a8161cc2153597ff7791f67e6cfde6a96ca2748292662ec78d0a"),
				},
				NodeAux: &proof.LeafNode{
					Key:   hashFromTextInt("15352856648520921629"),
					Value: hashFromInt(big.NewInt(0)),
				},
			},
		},
		{
			title: "test un-existence without aux node",
			//nolint:gocritic
			revNonce:    31, // revNonceKey[0] = 0b00011111
			revTreeRoot: hashFromInt(bigMerkleTree.Root().BigInt()),
			wantProof: proof.Proof{
				Existence: false,
				Siblings: []proof.Hash{
					hashFromHex("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					hashFromHex("28e5cdd29d9ad96cc214c654ca8e2f4fa5576bc132e172519804a58ee4bb4d18"),
					hashFromHex("658c7a65594ebb0815e1cc20f54284ccdb51bb1625f103c116ce58444145381e"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("5aa678402ef2cd5102de99722a6923183461b93f705a9d0aaaaff6a131a83504"),
				},
				NodeAux: nil,
			},
		},
		{
			title:       "test node does not exists",
			revNonce:    31,
			revTreeRoot: hashFromHex("1234567812345678123456781234567812345678123456781234567812345678"),
			wantErr:     "node not found",
		},
		{
			title:       "test zero tree root",
			revNonce:    31,
			revTreeRoot: hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
			wantProof: proof.Proof{
				Existence: false,
				Siblings:  nil,
				NodeAux:   nil,
			},
		},
		{
			title:       "existence of one only node in a tree",
			revNonce:    5577006791947779410,
			revTreeRoot: hashFromInt(oneNodeMerkleTree.Root().BigInt()),
			wantProof: proof.Proof{
				Existence: true,
				Siblings:  nil,
				NodeAux:   nil,
			},
		},
		{
			title:       "un-existence of one only node in a tree",
			revNonce:    10667007354186551956,
			revTreeRoot: hashFromInt(oneNodeMerkleTree.Root().BigInt()),
			wantProof: proof.Proof{
				Existence: false,
				Siblings:  nil,
				NodeAux: &proof.LeafNode{
					Key:   hashFromInt(big.NewInt(5577006791947779410)),
					Value: hashFromInt(big.NewInt(0)),
				},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			revNonceKey := hashFromInt(new(big.Int).SetUint64(tc.revNonce))
			revNonceValue := hashFromInt(big.NewInt(0))

			proofGen, err := proof.GenerateProof(rhsUrl, tc.revTreeRoot,
				revNonceKey)
			if tc.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.wantProof, proofGen)

				rootHash, err := proofGen.Root(revNonceKey, revNonceValue)
				require.NoError(t, err)
				require.Equal(t, tc.revTreeRoot, rootHash)

				//nolint:gocritic
				// logProof(t, proof)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
		})
	}

}

func getRevTreeRoot(rhsURL string, state proof.Hash) (proof.Hash, error) {
	stateNode, err := getNodeFromRHS(rhsURL, state)
	if err != nil {
		return proof.Hash{}, err
	}

	if len(stateNode.Children) != 3 {
		return proof.Hash{}, errors.New(
			"state hash does not looks like a state node: " +
				"number of children expected to be three")
	}

	return stateNode.Children[1], nil
}

var ErrNodeNotFound = errors.New("node not found")

func getNodeFromRHS(rhsURL string, hash proof.Hash) (proof.Node, error) {
	rhsURL = strings.TrimSuffix(rhsURL, "/")
	rhsURL += "/node/" + hash.Hex()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, rhsURL, http.NoBody)
	if err != nil {
		return proof.Node{}, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return proof.Node{}, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode == http.StatusNotFound {
		var resp map[string]interface{}
		dec := json.NewDecoder(httpResp.Body)
		err := dec.Decode(&resp)
		if err != nil {
			return proof.Node{}, err
		}
		if resp["status"] == "not found" {
			return proof.Node{}, ErrNodeNotFound
		} else {
			return proof.Node{}, errors.New("unexpected response")
		}
	} else if httpResp.StatusCode != http.StatusOK {
		return proof.Node{}, fmt.Errorf("unexpected response: %v",
			httpResp.StatusCode)
	}

	var nodeResp struct {
		Node   proof.Node `json:"node"`
		Status string     `json:"status"`
	}
	dec := json.NewDecoder(httpResp.Body)
	err = dec.Decode(&nodeResp)
	if err != nil {
		return proof.Node{}, err
	}

	return nodeResp.Node, nil
}

func saveIdenStateToRHS(t testing.TB, url string,
	merkleTree *merkletree.MerkleTree) proof.Hash {

	revTreeRoot := merkleTree.Root()
	state, err := poseidon.Hash([]*big.Int{big.NewInt(0), revTreeRoot.BigInt(),
		big.NewInt(0)})
	require.NoError(t, err)

	req := []proof.Node{
		{
			Hash: hashFromInt(state),
			Children: []proof.Hash{
				hashFromInt(merkletree.HashZero.BigInt()),
				hashFromInt(revTreeRoot.BigInt()),
				hashFromInt(merkletree.HashZero.BigInt())},
		},
	}
	submitNodesToRHS(t, url, req)
	return hashFromInt(state)
}

func buildTree(t testing.TB, revNonces []uint64) *merkletree.MerkleTree {
	mtStorage := memory.NewMemoryStorage()
	ctx := context.Background()
	const mtDepth = 40
	mt, err := merkletree.NewMerkleTree(ctx, mtStorage, mtDepth)
	require.NoError(t, err)

	for _, revNonce := range revNonces {
		key := new(big.Int).SetUint64(revNonce)
		value := big.NewInt(0)

		err = mt.Add(ctx, key, value)
		require.NoError(t, err)
	}

	return mt
}

func submitNodesToRHS(t testing.TB, url string, req []proof.Node) {
	reqBytes, err := json.Marshal(req)
	require.NoError(t, err)
	bodyReader := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequest(http.MethodPost, url+"/node", bodyReader)
	require.NoError(t, err)

	httpResp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, httpResp.Body.Close())
	}()
	respData, err := io.ReadAll(httpResp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, httpResp.StatusCode)
	require.Equal(t, `{"status":"OK"}`, string(respData))
}

func saveTreeToRHS(t testing.TB, url string,
	merkleTree *merkletree.MerkleTree) {
	ctx := context.Background()
	var req []proof.Node
	hashOne := merkletree.NewHashFromBigInt(big.NewInt(1))
	err := merkleTree.Walk(ctx, nil, func(node *merkletree.Node) {
		nodeKey, err := node.Key()
		require.NoError(t, err)
		switch node.Type {
		case merkletree.NodeTypeMiddle:
			req = append(req, proof.Node{
				Hash: hashFromHex(nodeKey.Hex()),
				Children: []proof.Hash{
					hashFromHex(node.ChildL.Hex()),
					hashFromHex(node.ChildR.Hex())},
			})
		case merkletree.NodeTypeLeaf:
			req = append(req, proof.Node{
				Hash: hashFromHex(nodeKey.Hex()),
				Children: []proof.Hash{
					hashFromHex(node.Entry[0].Hex()),
					hashFromHex(node.Entry[1].Hex()),
					hashFromHex(hashOne.Hex())},
			})
		case merkletree.NodeTypeEmpty:
			// do not save zero nodes
		default:
			require.Failf(t, "unexpected node type", "unexpected node type: %v",
				node.Type)
		}
	})
	require.NoError(t, err)

	submitNodesToRHS(t, url, req)
}

func hashFromHex(in string) proof.Hash {
	h, err := proof.NewHashFromHex(in)
	if err != nil {
		panic(err)
	}
	return h
}

func hashFromInt(in *big.Int) proof.Hash {
	h, err := proof.NewHashFromBigInt(in)
	if err != nil {
		panic(err)
	}
	return h
}

func hashFromTextInt(in string) proof.Hash {
	i, ok := new(big.Int).SetString(in, 10)
	if !ok {
		panic(in)
	}
	return hashFromInt(i)
}
