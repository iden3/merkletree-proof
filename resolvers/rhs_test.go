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
	client, err := ethclient.Dial("")
	require.NoError(t, err)
	var ethClients map[core.ChainID]*ethclient.Client = make(map[core.ChainID]*ethclient.Client)
	ethClients[80001] = client

	opts := []RHSResolverOpts{WithEthClients(ethClients), WithIssuerDID(issuerDID), WithStateContractAddr(stateAddr)}

	config := RHSResolverConfig{}
	for _, o := range opts {
		o(&config)
	}

	rhsResolver := RHSResolver{config}
	_, err = rhsResolver.Resolve(context.Background(), credStatus)
	require.NoError(t, err)

}
