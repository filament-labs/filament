package walletclient

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

// WalletClient wraps a Filecoin full node API connection.
type WalletClient struct {
	api    api.FullNode
	closer jsonrpc.ClientCloser
	url    string
}

// New connects to the Lotus RPC endpoint and returns a reusable WalletClient.
func New(rpcURL string) (*WalletClient, error) {
	var apiNode api.FullNode

	closer, err := jsonrpc.NewClient(context.TODO(), rpcURL, "Filecoin", &apiNode, nil)
	if err != nil {
		return nil, err
	}

	return &WalletClient{
		api:    apiNode,
		closer: closer,
		url:    rpcURL,
	}, nil
}

// Close releases resources held by the underlying RPC client.
// Should be called once when shutting down the app or service.
func (wc *WalletClient) Close() {
	wc.closer()
}

func (wc *WalletClient) SendTransaction(to string, amount types.BigInt, passphrase string) {

}

func (wc *WalletClient) WalletList(ctx context.Context) ([]string, error) {
	return []string{}, nil
}
