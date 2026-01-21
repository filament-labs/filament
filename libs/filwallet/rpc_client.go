package filwallet

import (
	"context"
	"fmt"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
)

type RPCClient struct {
	node   api.FullNode
	closer jsonrpc.ClientCloser
}

func NewRPCClient(ctx context.Context, rpcEndpoint, rpcToken string) (*RPCClient, error) {
	headers := make(http.Header)
	if rpcToken != "" {
		headers.Set("Authorization", "Bearer "+rpcToken)
	}

	node, closer, err := client.NewFullNodeRPCV1(ctx, rpcEndpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("dial rpc %s: %w", rpcEndpoint, err)
	}

	return &RPCClient{
		node:   node,
		closer: closer,
	}, nil
}

func (c *RPCClient) Close() {
	if c.closer == nil {
		return
	}

	c.closer()
	c.closer = nil
}
