package merkletree_proof

import (
	"encoding/json"
	"testing"

	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/stretchr/testify/require"
)

func TestNode_MarshalJSON(t *testing.T) {
	testCases := []struct {
		title string
		in    Node
		want  string
	}{
		{
			title: "regular node",
			in: Node{
				Hash: hashFromHex("20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222"),
				Children: []*merkletree.Hash{
					hashFromHex("79f66791900bc0c9260f708e317437415b9f45673384f5b0752f5a649f661207"),
					hashFromHex("f9b198c1da06c8cc8aedf408f2be2fd9def1818496924542c3194ceb7c70bb01"),
					hashFromHex("4012c3753476058e08d36af518ba61ea65b49c0318af0bd976c95a931e257b28"),
				},
			},
			want: `{
  "hash": "20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222",
  "children": [
	"79f66791900bc0c9260f708e317437415b9f45673384f5b0752f5a649f661207",
	"f9b198c1da06c8cc8aedf408f2be2fd9def1818496924542c3194ceb7c70bb01",
	"4012c3753476058e08d36af518ba61ea65b49c0318af0bd976c95a931e257b28"
  ]
}`,
		},
		{
			title: "empty children",
			in: Node{
				Hash:     hashFromHex("20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222"),
				Children: []*merkletree.Hash{},
			},
			want: `{
  "hash": "20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222",
  "children": []
}`,
		},
		{
			title: "nil children",
			in: Node{
				Hash:     hashFromHex("20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222"),
				Children: nil,
			},
			want: `{
  "hash": "20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222",
  "children": null
}`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			result, err := json.Marshal(tc.in)
			require.NoError(t, err)
			require.JSONEq(t, tc.want, string(result))
		})
	}
}

func TestNode_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		title string
		in    string
		want  Node
	}{
		{
			title: "regular node",
			in: `{
  "hash": "20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222",
  "children": [
	"79f66791900bc0c9260f708e317437415b9f45673384f5b0752f5a649f661207",
	"f9b198c1da06c8cc8aedf408f2be2fd9def1818496924542c3194ceb7c70bb01",
	"4012c3753476058e08d36af518ba61ea65b49c0318af0bd976c95a931e257b28"
  ]
}`,
			want: Node{
				Hash: hashFromHex("20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222"),
				Children: []*merkletree.Hash{
					hashFromHex("79f66791900bc0c9260f708e317437415b9f45673384f5b0752f5a649f661207"),
					hashFromHex("f9b198c1da06c8cc8aedf408f2be2fd9def1818496924542c3194ceb7c70bb01"),
					hashFromHex("4012c3753476058e08d36af518ba61ea65b49c0318af0bd976c95a931e257b28"),
				},
			},
		},
		{
			title: "empty children",
			in: `{
  "hash": "20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222",
  "children": []
}`,
			want: Node{
				Hash:     hashFromHex("20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222"),
				Children: []*merkletree.Hash{},
			},
		},
		{
			title: "nil children",
			in: `{
  "hash": "20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222",
  "children": null
}`,
			want: Node{
				Hash:     hashFromHex("20a8bc6b66482191ad30d7c0a95e7a512297f0a2da9fccc0803b0b03aa3f5222"),
				Children: nil,
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			var n Node
			err := json.Unmarshal([]byte(tc.in), &n)
			require.NoError(t, err)
			require.Equal(t, tc.want, n)
		})
	}
}

func hashFromHex(in string) *merkletree.Hash {
	h, err := merkletree.NewHashFromHex(in)
	if err != nil {
		panic(err)
	}
	return h
}
