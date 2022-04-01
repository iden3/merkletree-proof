package merkletree_proof

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-crypto/constants"
	"github.com/stretchr/testify/require"
)

func TestNewHashFromBigInt(t *testing.T) {
	testCases := []struct {
		title   string
		input   *big.Int
		want    Hash
		wantErr string
	}{
		{
			title: "zero",
			input: big.NewInt(0),
			want:  Hash{},
		},
		{
			title: "one",
			input: big.NewInt(1),
			want:  Hash{0x01},
		},
		{
			title: "max",
			input: new(big.Int).Sub(constants.Q, big.NewInt(1)),
			want: Hash{
				0x0, 0x0, 0x0, 0xf0, 0x93, 0xf5, 0xe1, 0x43,
				0x91, 0x70, 0xb9, 0x79, 0x48, 0xe8, 0x33, 0x28,
				0x5d, 0x58, 0x81, 0x81, 0xb6, 0x45, 0x50, 0xb8,
				0x29, 0xa0, 0x31, 0xe1, 0x72, 0x4e, 0x64, 0x30},
		},
		{
			title:   "too large",
			input:   constants.Q,
			wantErr: "big int out of field",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			got, err := NewHashFromBigInt(tc.input)
			if tc.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
		})
	}
}

func TestHash_Hex(t *testing.T) {
	in := Hash{
		0x0, 0x0, 0x0, 0xf0, 0x93, 0xf5, 0xe1, 0x43,
		0x91, 0x70, 0xb9, 0x79, 0x48, 0xe8, 0x33, 0x28,
		0x5d, 0x58, 0x81, 0x81, 0xb6, 0x45, 0x50, 0xb8,
		0x29, 0xa0, 0x31, 0xe1, 0x72, 0x4e, 0x64, 0x30}
	want := "000000f093f5e1439170b97948e833285d588181b64550b829a031e1724e6430"
	require.Equal(t, want, in.Hex())
}

func TestHash_Int(t *testing.T) {
	in := Hash{
		0x0, 0x0, 0x0, 0xf0, 0x93, 0xf5, 0xe1, 0x43,
		0x91, 0x70, 0xb9, 0x79, 0x48, 0xe8, 0x33, 0x28,
		0x5d, 0x58, 0x81, 0x81, 0xb6, 0x45, 0x50, 0xb8,
		0x29, 0xa0, 0x31, 0xe1, 0x72, 0x4e, 0x64, 0x30}
	want := new(big.Int).Sub(constants.Q, big.NewInt(1))
	require.Equal(t, 0, in.Int().Cmp(want))
}

func TestProof_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		title   string
		in      string
		want    Proof
		wantErr string
	}{
		{
			title: "OK",
			in: `{
"existence": true,
"siblings": null}`,
			want: Proof{Existence: true},
		},
		{
			title: "only existence",
			in:    `{"existence": true}`,
			want:  Proof{Existence: true},
		},
		{
			title: "null siblings",
			in: `{
  "existence": true,
  "siblings": null
}`,
			want: Proof{Existence: true},
		},
		{
			title: "empty siblings",
			in: `{
  "existence": true,
  "siblings": []
}`,
			want: Proof{Existence: true, Siblings: make([]Hash, 0)},
		},
		{
			title: "with siblings",
			in: `{
  "existence": true,
  "siblings": [
    "b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14",
    "74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"
  ]
}`,
			want: Proof{
				Existence: true,
				Siblings: []Hash{
					mkHash("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					mkHash("74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"),
				}},
		},
		{
			title: "sibling out of field",
			in: `{
  "existence": true,
  "siblings": [
    "b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14",
    "020000f093f5e1439170b97948e833285d588181b64550b829a031e1724e6430"
  ]
}`,
			wantErr: "big int out of field",
		},
		{
			title: "with aux_node",
			in: `{
  "existence": true,
  "siblings": [
    "b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14",
    "74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"
  ],
  "aux_node": {
    "key":   "94d2c422acd20894000000000000000000000000000000000000000000000000",
    "value": "0000000000000000000000000000000000000000000000000000000000000000"
  }
}`,
			want: Proof{
				Existence: true,
				Siblings: []Hash{
					mkHash("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					mkHash("74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"),
				},
				NodeAux: &LeafNode{
					Key:   mkHash("94d2c422acd20894000000000000000000000000000000000000000000000000"),
					Value: mkHash("0000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			var p Proof
			err := json.Unmarshal([]byte(tc.in), &p)
			if tc.wantErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.want, p)
			} else {
				require.EqualError(t, err, tc.wantErr)
			}
		})
	}
}

func mkHash(in string) Hash {
	h, err := NewHashFromHex(in)
	if err != nil {
		panic(err)
	}
	return h
}

func TestProof_MarshalJSON(t *testing.T) {
	testCases := []struct {
		title string
		in    Proof
		want  string
	}{
		{
			title: "all fields filled",
			in: Proof{
				Existence: true,
				Siblings: []Hash{
					mkHash("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					mkHash("74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"),
				},
				NodeAux: &LeafNode{
					Key:   mkHash("94d2c422acd20894000000000000000000000000000000000000000000000000"),
					Value: mkHash("0000000000000000000000000000000000000000000000000000000000000000"),
				},
			},
			want: `{
  "existence": true,
  "siblings": [
    "b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14",
    "74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"
  ],
  "aux_node": {
    "key":   "94d2c422acd20894000000000000000000000000000000000000000000000000",
    "value": "0000000000000000000000000000000000000000000000000000000000000000"
  }
}`,
		},
		{
			title: "without aux_node",
			in: Proof{
				Existence: true,
				Siblings: []Hash{
					mkHash("b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14"),
					mkHash("74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"),
				},
				NodeAux: nil,
			},
			want: `{
  "existence": true,
  "siblings": [
    "b2f5a640931d3815375be1e9a00ee4da175d3eb9520ef0715f484b11a75f2a14",
    "74321998e281c0a89dbcce55a6cec0e366536e2697ea40efaf036ecba751ed03"
  ]
}`,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			res, err := json.Marshal(tc.in)
			require.NoError(t, err)
			require.JSONEq(t, tc.want, string(res))
		})
	}
}
