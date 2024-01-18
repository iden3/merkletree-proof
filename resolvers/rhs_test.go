package resolvers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	core "github.com/iden3/go-iden3-core/v2"
	"github.com/iden3/go-iden3-core/v2/w3c"
	"github.com/iden3/go-schema-processor/v2/verifiable"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestRhsResolver(t *testing.T) {
	credStatusJSON := `{
		"id": "https://rhs-staging.polygonid.me/node?state=f9dd6aa4e1abef52b6c94ab7eb92faf1a283b371d263e25ac835c9c04894741e",
		"revocationNonce": 0,
		"statusIssuer": {
			"id": "https://ad40-91-210-251-7.ngrok-free.app/api/v1/identities/did%3Apolygonid%3Apolygon%3Amumbai%3A2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf/claims/revocation/status/0",
			"revocationNonce": 0,
			"type": "SparseMerkleTreeProof"
		},
		"type": "Iden3ReverseSparseMerkleTreeProof"
	}`

	var credStatus verifiable.CredentialStatus
	err := json.Unmarshal([]byte(credStatusJSON), &credStatus)
	require.NoError(t, err)

	issuerDID, err := w3c.ParseDID("did:polygonid:polygon:mumbai:2qLGnFZiHrhdNh5KwdkGvbCN1sR2pUaBpBahAXC3zf")
	require.NoError(t, err)
	stateAddr := common.HexToAddress("0x134B1BE34911E39A8397ec6289782989729807a4")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterMatcherResponder("POST", "http://my-rpc/v2/1111",
		httpmock.BodyContainsString(`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"data":"0xb4bdea550010961e749448c0c935c85ae263d271b383a2f1fa92ebb74ac9b652efab1202","from":"0x0000000000000000000000000000000000000000","to":"0x134b1be34911e39a8397ec6289782989729807a4"},"latest"]}`),
		httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":1,"result":"0x0010961e749448c0c935c85ae263d271b383a2f1fa92ebb74ac9b652efab120209444d55300819a07594d731c0bf7d3713f5e9324e0435f926c3ef1d8e4a823400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000065846207000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029cf4ff0000000000000000000000000000000000000000000000000000000000000000"}`))

	httpmock.RegisterResponder("GET", "https://rhs-staging.polygonid.me/node/34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409",
		httpmock.NewStringResponder(200, `{"node":{"hash":"34824a8e1defc326f935044e32e9f513377dbfc031d79475a0190830554d4409","children":["4436ea12d352ddb84d2ac7a27bbf7c9f1bfc7d3ff69f3e6cf4348f424317fd0b","0000000000000000000000000000000000000000000000000000000000000000","37eabc712cdaa64793561b16b8143f56f149ad1b0c35297a1b125c765d1c071e"]},"status":"OK"}`))

	client, err := ethclient.Dial("http://my-rpc/v2/1111")
	require.NoError(t, err)
	var ethClients map[core.ChainID]*ethclient.Client = make(map[core.ChainID]*ethclient.Client)
	ethClients[80001] = client

	config := RHSResolverConfig{
		EthClients:        ethClients,
		StateContractAddr: stateAddr,
	}

	rhsResolver := NewRHSResolver(config)
	_, err = rhsResolver.Resolve(context.Background(), credStatus, &verifiable.CredentialStatusResolveOptions{IssuerDID: issuerDID})
	require.NoError(t, err)

}
