package integration

import (
	"context"
	"errors"
	"math/big"
	"os"
	//"os"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"github.com/iden3/merkletree-proof/common"
	proof "github.com/iden3/merkletree-proof/http"
	"github.com/stretchr/testify/require"
)

func TestProof_Http(t *testing.T) {
	t.Skip("skipping http test")
	rhsUrl, ok := os.LookupEnv("RHS_URL")
	if !ok || rhsUrl == "" {
		t.Fatal("RHS_URL not set")
	}
	rhsCli := &proof.HTTPReverseHashCli{URL: rhsUrl}

	runTestCases(t, rhsCli)
}

func TestProof_Eth(t *testing.T) {
	signer := NewTestSigner()

	addrStr, ok := os.LookupEnv("IDENTITY_TREE_STORE_ADDRESS")
	if !ok {
		panic("IDENTITY_TREE_STORE_ADDRESS not set")
	}

	addr := ethcommon.HexToAddress(addrStr)

	cli, err := NewTestEthRpcReserveHashCli(addr, signer)
	if err != nil {
		panic(err)
	}

	runTestCases(t, cli)
}

func runTestCases(t *testing.T, rhsCli common.ReverseHashCli) {
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
	saveTreeToRHS(t, rhsCli, bigMerkleTree)

	oneNodeMerkleTree := buildTree(t, []uint64{5577006791947779410})
	saveTreeToRHS(t, rhsCli, oneNodeMerkleTree)

	t.Run("Test save state", func(t *testing.T) {
		state := saveIdenStateToRHS(t, rhsCli, bigMerkleTree)

		revTreeRoot, err := getRevTreeRoot(rhsCli, state)
		require.NoError(t, err)

		require.Equal(t, bigMerkleTree.Root(), revTreeRoot)
	})

	testCases := []struct {
		title       string
		revNonce    uint64
		revTreeRoot *merkletree.Hash
		wantProof   *merkletree.Proof
		wantErr     string
	}{
		{
			title:       "regular node",
			revNonce:    10667007354186551956,
			revTreeRoot: bigMerkleTree.Root(),
			wantProof: mkProof(
				true,
				[]*merkletree.Hash{
					hashFromHex("74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"),
					hashFromHex("ff11b8bf1d13e28e86e249d2acdba0bd9c0fe4a5f56ad4236b09185bde81c316"),
					hashFromHex("db5eb80f6b60b4e23714d4d00f178ba62fbdb4f0294675f51ac99aa24e600827"),
				},
				nil),
		},
		{
			title:       "a node with zero siblings",
			revNonce:    8674665223082147919,
			revTreeRoot: bigMerkleTree.Root(),
			wantProof: mkProof(
				true,
				[]*merkletree.Hash{
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
				nil),
		},
		{
			title: "un-existence with aux node",
			//nolint:gocritic
			revNonce:    5, // revNonceKey[0] = 0b00000101
			revTreeRoot: bigMerkleTree.Root(),
			wantProof: mkProof(
				false,
				[]*merkletree.Hash{
					hashFromHex("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					hashFromHex("c9719432e3d8bf360d0f2de456c5321c51295895c9330b0588552580765cd929"),
					hashFromHex("c0e8bf477403a8161cc2153597ff7791f67e6cfde6a96ca2748292662ec78d0a"),
				},
				&merkletree.NodeAux{
					Key: hashFromInt(
						new(big.Int).SetUint64(15352856648520921629)),
					Value: &merkletree.HashZero,
				}),
		},
		{
			title: "test un-existence without aux node",
			//nolint:gocritic
			revNonce:    31, // revNonceKey[0] = 0b00011111
			revTreeRoot: bigMerkleTree.Root(),
			wantProof: mkProof(
				false,
				[]*merkletree.Hash{
					hashFromHex("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					hashFromHex("28e5cdd29d9ad96cc214c654ca8e2f4fa5576bc132e172519804a58ee4bb4d18"),
					hashFromHex("658c7a65594ebb0815e1cc20f54284ccdb51bb1625f103c116ce58444145381e"),
					hashFromHex("0000000000000000000000000000000000000000000000000000000000000000"),
					hashFromHex("5aa678402ef2cd5102de99722a6923183461b93f705a9d0aaaaff6a131a83504"),
				},
				nil,
			),
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
			wantProof:   mkProof(false, nil, nil),
		},
		{
			title:       "existence of one only node in a tree",
			revNonce:    5577006791947779410,
			revTreeRoot: oneNodeMerkleTree.Root(),
			wantProof:   mkProof(true, nil, nil),
		},
		{
			title:       "un-existence of one only node in a tree",
			revNonce:    10667007354186551956,
			revTreeRoot: oneNodeMerkleTree.Root(),
			wantProof: mkProof(false, nil,
				&merkletree.NodeAux{
					Key: hashFromInt(
						big.NewInt(5577006791947779410)),
					Value: &merkletree.HashZero,
				},
			),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			revNonceKeyInt := new(big.Int).SetUint64(tc.revNonce)
			revNonceKey := hashFromInt(revNonceKeyInt)
			revNonceValueInt := big.NewInt(0)

			proofGen, err := rhsCli.GenerateProof(context.Background(),
				tc.revTreeRoot, revNonceKey)
			if tc.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.wantProof, proofGen)

				rootHash, err := merkletree.RootFromProof(proofGen,
					revNonceKeyInt, revNonceValueInt)
				require.NoError(t, err)
				require.Equal(t, tc.revTreeRoot, rootHash)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
		})
	}
}

func getRevTreeRoot(rhsCli common.ReverseHashCli,
	state *merkletree.Hash) (*merkletree.Hash, error) {
	stateNode, err := rhsCli.GetNode(context.Background(), state)
	if err != nil {
		return nil, err
	}

	if len(stateNode.Children) != 3 {
		return nil, errors.New(
			"state hash does not looks like a state node: " +
				"number of children expected to be three")
	}

	return stateNode.Children[1], nil
}

func saveIdenStateToRHS(t testing.TB, rhsCli common.ReverseHashCli,
	merkleTree *merkletree.MerkleTree) *merkletree.Hash {

	revTreeRoot := merkleTree.Root()
	state, err := poseidon.Hash([]*big.Int{big.NewInt(0), revTreeRoot.BigInt(),
		big.NewInt(0)})
	require.NoError(t, err)

	stateHash := hashFromInt(state)
	req := []common.Node{
		{
			Hash: stateHash,
			Children: []*merkletree.Hash{&merkletree.HashZero,
				revTreeRoot, &merkletree.HashZero},
		},
	}
	err = rhsCli.SaveNodes(context.Background(), req)
	require.NoError(t, err)
	return stateHash
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

func saveTreeToRHS(t testing.TB, rhsCli common.ReverseHashCli,
	merkleTree *merkletree.MerkleTree) {
	ctx := context.Background()
	var req []common.Node
	hashOne := hashFromInt(big.NewInt(1))
	err := merkleTree.Walk(ctx, nil, func(node *merkletree.Node) {
		nodeKey, err := node.Key()
		require.NoError(t, err)
		switch node.Type {
		case merkletree.NodeTypeMiddle:
			req = append(req, common.Node{
				Hash:     nodeKey,
				Children: []*merkletree.Hash{node.ChildL, node.ChildR}})
		case merkletree.NodeTypeLeaf:
			req = append(req, common.Node{
				Hash: nodeKey,
				Children: []*merkletree.Hash{node.Entry[0], node.Entry[1],
					hashOne},
			})
		case merkletree.NodeTypeEmpty:
			// do not save zero nodes
		default:
			require.Failf(t, "unexpected node type", "unexpected node type: %v",
				node.Type)
		}
	})
	require.NoError(t, err)

	err = rhsCli.SaveNodes(context.Background(), req)
	require.NoError(t, err)
}

func hashFromHex(in string) *merkletree.Hash {
	h, err := merkletree.NewHashFromHex(in)
	if err != nil {
		panic(err)
	}
	return h
}

func mkProof(existence bool, siblings []*merkletree.Hash,
	nodeAux *merkletree.NodeAux) *merkletree.Proof {
	p, err := merkletree.NewProofFromData(existence, siblings, nodeAux)
	if err != nil {
		panic(err)
	}
	return p
}

func hashFromInt(in *big.Int) *merkletree.Hash {
	h, err := merkletree.NewHashFromBigInt(in)
	if err != nil {
		panic(err)
	}
	return h
}
