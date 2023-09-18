package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/merkletree-proof/common"
)

func init() {
	hashOneP, err := merkletree.NewHashFromBigInt(big.NewInt(1))
	if err != nil {
		panic(err)
	}
	copy(hashOne[:], hashOneP[:])
}

var hashOne merkletree.Hash

var ErrNodeNotFound = errors.New("node not found")

type HTTPReverseHashCli struct {
	URL         string
	HTTPTimeout time.Duration
}

// GenerateProof generates proof of existence or in-existence of a key in
// a tree identified by a treeRoot.
func (cli *HTTPReverseHashCli) GenerateProof(ctx context.Context,
	treeRoot *merkletree.Hash,
	key *merkletree.Hash) (*merkletree.Proof, error) {

	if cli.URL == "" {
		return nil, errors.New("HTTP reverse hash service url is not specified")
	}

	return common.GenerateProof(ctx, cli, treeRoot, key)
}

func (cli *HTTPReverseHashCli) nodeURL(node *merkletree.Hash) string {
	nodeURL := cli.baseURL() + "/node"
	if node == nil {
		return nodeURL
	}
	return nodeURL + "/" + node.Hex()
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

func (cli *HTTPReverseHashCli) GetNode(ctx context.Context,
	hash *merkletree.Hash) (common.Node, error) {

	if hash == nil {
		return common.Node{}, errors.New("hash is nil")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cli.getHttpTimeout())
		defer cancel()
	}

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodGet, cli.nodeURL(hash), http.NoBody)
	if err != nil {
		return common.Node{}, err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return common.Node{}, err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode == http.StatusNotFound {
		var resp map[string]interface{}
		dec := json.NewDecoder(httpResp.Body)
		err := dec.Decode(&resp)
		if err != nil {
			return common.Node{}, err
		}
		if resp["status"] == "not found" {
			return common.Node{}, ErrNodeNotFound
		} else {
			return common.Node{}, errors.New("unexpected response")
		}
	} else if httpResp.StatusCode != http.StatusOK {
		return common.Node{}, fmt.Errorf("unexpected response: %v",
			httpResp.StatusCode)
	}

	var nodeResp nodeResponse
	dec := json.NewDecoder(httpResp.Body)
	err = dec.Decode(&nodeResp)
	if err != nil {
		return common.Node{}, err
	}

	return nodeResp.Node, nil
}

func (cli *HTTPReverseHashCli) SaveNodes(ctx context.Context,
	nodes []common.Node) error {

	reqBytes, err := json.Marshal(nodes)
	if err != nil {
		return err
	}

	// if no timeout set on context, set it here
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cli.getHttpTimeout())
		defer cancel()
	}

	bodyReader := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		cli.nodeURL(nil), bodyReader)
	if err != nil {
		return err
	}

	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	dec := json.NewDecoder(httpResp.Body)
	var respM map[string]interface{}
	err = dec.Decode(&respM)
	if err != nil {
		return fmt.Errorf("unable to decode RHS response: %w", err)
	}

	if respM["status"] != "OK" {
		return fmt.Errorf("unexpected RHS response status: %s", respM["status"])
	}

	return nil
}

type nodeResponse struct {
	Node   common.Node `json:"node"`
	Status string      `json:"status"`
}
