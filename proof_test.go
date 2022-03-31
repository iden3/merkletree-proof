package merkletree_proof

import (
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
			want: Hash{0x0, 0x0, 0x0, 0xf0, 0x93, 0xf5, 0xe1, 0x43, 0x91, 0x70,
				0xb9, 0x79, 0x48, 0xe8, 0x33, 0x28, 0x5d, 0x58, 0x81, 0x81,
				0xb6, 0x45, 0x50, 0xb8, 0x29, 0xa0, 0x31, 0xe1, 0x72, 0x4e,
				0x64, 0x30},
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
