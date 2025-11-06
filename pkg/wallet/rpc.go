package wallet

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
	cid "github.com/ipfs/go-cid"
)

// RPCClient wraps Filecoin RPC API interactions with clean lifecycle management.
type RPCClient struct {
	api    api.FullNode
	closer jsonrpc.ClientCloser
	config RPCConfig
}

// NewRPCClient creates a new RPC client with validated config and authenticated connection.
func NewRPCClient(config RPCConfig) (*RPCClient, error) {
	if config.Endpoint == "" {
		return nil, fmt.Errorf("RPC endpoint is required")
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	headers := make(map[string][]string)
	if config.Token != "" {
		headers["Authorization"] = []string{"Bearer " + config.Token}
	}

	var apiNode api.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(
		ctx,
		config.Endpoint,
		"Filecoin",
		[]interface{}{
			&apiNode.Internal,
			&apiNode.CommonStruct.Internal,
		},
		headers,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Filecoin RPC: %w", err)
	}

	return &RPCClient{
		api:    &apiNode,
		closer: closer,
		config: config,
	}, nil
}

// Close safely closes the RPC connection.
func (c *RPCClient) Close() {
	if c.closer != nil {
		c.closer()
	}
}

// GetBalance returns the current balance and nonce for an address.
func (c *RPCClient) GetBalance(ctx context.Context, addr address.Address) (*Balance, error) {
	actor, err := c.api.StateGetActor(ctx, addr, types.EmptyTSK)
	if err != nil {
		return nil, fmt.Errorf("failed to get actor state: %w", err)
	}

	return &Balance{
		Address:   addr.String(),
		Balance:   actor.Balance.String(),
		Nonce:     actor.Nonce,
		Timestamp: time.Now().UTC(),
	}, nil
}

// GetNonce fetches the current nonce for the given address.
func (c *RPCClient) GetNonce(ctx context.Context, addr address.Address) (uint64, error) {
	actor, err := c.api.StateGetActor(ctx, addr, types.EmptyTSK)
	if err != nil {
		return 0, fmt.Errorf("failed to get actor for nonce: %w", err)
	}
	return actor.Nonce, nil
}

// EstimateGas computes gas parameters for a message.
func (c *RPCClient) EstimateGas(ctx context.Context, msg *types.Message) (*GasEstimate, error) {
	limit, err := c.api.GasEstimateGasLimit(ctx, msg, types.EmptyTSK)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas limit: %w", err)
	}

	gasPremium, err := c.api.GasEstimateGasPremium(ctx, 10, msg.From, limit, types.EmptyTSK)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas premium: %w", err)
	}

	feeCap, err := c.api.GasEstimateFeeCap(ctx, msg, 20, types.EmptyTSK)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate fee cap: %w", err)
	}

	return &GasEstimate{
		GasLimit:   limit,
		GasFeeCap:  feeCap.String(),
		GasPremium: gasPremium.String(),
	}, nil
}

// SendMessage pushes a signed message to the mempool and returns its CID.
func (c *RPCClient) SendMessage(ctx context.Context, signedMsg *types.SignedMessage) (string, error) {
	cid, err := c.api.MpoolPush(ctx, signedMsg)
	if err != nil {
		return "", fmt.Errorf("failed to push message to mempool: %w", err)
	}
	return cid.String(), nil
}

// WaitForMessage waits for a message to be included in a tipset with given confidence.
func (c *RPCClient) WaitForMessage(ctx context.Context, cidStr string, confidence uint64) (*types.MessageReceipt, error) {
	msgCid, err := parseCID(cidStr)
	if err != nil {
		return nil, err
	}

	receipt, err := c.api.StateWaitMsg(ctx, msgCid, confidence, 0, true)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for message inclusion: %w", err)
	}

	return &receipt.Receipt, nil
}

// GetTransaction retrieves full transaction details by CID.
func (c *RPCClient) GetTransaction(ctx context.Context, cidStr string) (*Transaction, error) {
	msgCid, err := parseCID(cidStr)
	if err != nil {
		return nil, err
	}

	msgLookup, err := c.api.StateSearchMsg(ctx, types.EmptyTSK, msgCid, -1, false)
	if err != nil {
		return nil, fmt.Errorf("message not found in chain: %w", err)
	}
	if msgLookup == nil {
		return nil, fmt.Errorf("message CID %s not found", cidStr)
	}

	msg, err := c.api.ChainGetMessage(ctx, msgCid)
	if err != nil {
		return nil, fmt.Errorf("failed to get message content: %w", err)
	}

	ts, err := c.api.ChainGetTipSet(ctx, msgLookup.TipSet)
	if err != nil {
		return nil, fmt.Errorf("failed to get tipset: %w", err)
	}

	status := "confirmed"
	if msgLookup.Receipt.ExitCode != 0 {
		status = "failed"
	}

	blockTime := time.Unix(int64(ts.Blocks()[0].Timestamp), 0).UTC()

	return &Transaction{
		Cid:        cidStr,
		From:       msg.From.String(),
		To:         msg.To.String(),
		Value:      msg.Value.String(),
		GasFeeCap:  msg.GasFeeCap.String(),
		GasPremium: msg.GasPremium.String(),
		GasLimit:   msg.GasLimit,
		Nonce:      msg.Nonce,
		Method:     uint64(msg.Method),
		Params:     msg.Params,
		Timestamp:  blockTime,
		Height:     int64(msgLookup.Height),
		Status:     status,
	}, nil
}

// GetTransactions returns up to `limit` recent transactions involving the address.
func (c *RPCClient) GetTransactions(ctx context.Context, addr address.Address, limit int) ([]*Transaction, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	head, err := c.api.ChainHead(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain head: %w", err)
	}

	transactions := make([]*Transaction, 0, limit)
	currentHeight := head.Height()
	maxSearchDepth := abi.ChainEpoch(2000)

	for len(transactions) < limit && currentHeight > 0 {
		if head.Height()-currentHeight > maxSearchDepth {
			break
		}

		ts, err := c.api.ChainGetTipSetByHeight(ctx, currentHeight, types.EmptyTSK)
		if err != nil {
			currentHeight--
			continue
		}

		for _, block := range ts.Blocks() {
			msgs, err := c.api.ChainGetBlockMessages(ctx, block.Cid())
			if err != nil {
				continue
			}

			c.appendMatchingBLSMessages(&transactions, msgs.BlsMessages, addr, block.Timestamp, uint64(currentHeight), limit)
			c.appendMatchingSECPMessages(&transactions, msgs.SecpkMessages, addr, block.Timestamp, uint64(currentHeight), limit)

			if len(transactions) >= limit {
				return transactions, nil
			}
		}

		currentHeight--
	}

	return transactions, nil
}

// appendMatchingBLSMessages appends BLS messages involving addr to txs.
func (c *RPCClient) appendMatchingBLSMessages(
	txs *[]*Transaction,
	msgs []*types.Message,
	addr address.Address,
	timestamp uint64,
	height uint64,
	limit int,
) {
	for _, msg := range msgs {
		if msg.From == addr || msg.To == addr {
			*txs = append(*txs, &Transaction{
				From:       msg.From.String(),
				To:         msg.To.String(),
				Value:      msg.Value.String(),
				GasFeeCap:  msg.GasFeeCap.String(),
				GasPremium: msg.GasPremium.String(),
				GasLimit:   msg.GasLimit,
				Nonce:      msg.Nonce,
				Method:     uint64(msg.Method),
				Timestamp:  time.Unix(int64(timestamp), 0).UTC(),
				Height:     int64(height),
				Status:     "confirmed",
			})
			if len(*txs) >= limit {
				return
			}
		}
	}
}

// appendMatchingSECPMessages appends SECP messages involving addr to txs.
func (c *RPCClient) appendMatchingSECPMessages(
	txs *[]*Transaction,
	msgs []*types.SignedMessage,
	addr address.Address,
	timestamp uint64,
	height uint64,
	limit int,
) {
	for _, smsg := range msgs {
		msg := smsg.Message
		if msg.From == addr || msg.To == addr {
			*txs = append(*txs, &Transaction{
				From:       msg.From.String(),
				To:         msg.To.String(),
				Value:      msg.Value.String(),
				GasFeeCap:  msg.GasFeeCap.String(),
				GasPremium: msg.GasPremium.String(),
				GasLimit:   msg.GasLimit,
				Nonce:      msg.Nonce,
				Method:     uint64(msg.Method),
				Timestamp:  time.Unix(int64(timestamp), 0).UTC(),
				Height:     int64(height),
				Status:     "confirmed",
			})
			if len(*txs) >= limit {
				return
			}
		}
	}
}

// GetChainHead returns the current chain height.
func (c *RPCClient) GetChainHead(ctx context.Context) (int64, error) {
	head, err := c.api.ChainHead(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get chain head: %w", err)
	}
	return int64(head.Height()), nil
}

// parseCID safely converts a CID string into a cid.Cid.
func parseCID(cidStr string) (cid.Cid, error) {
	c, err := cid.Decode(cidStr)
	if err != nil {
		return cid.Undef, fmt.Errorf("invalid CID format: %w", err)
	}
	return c, nil
}

// ParseFIL converts a human-readable FIL string (e.g. "1.5") to attoFIL string.
func ParseFIL(fil string) (string, error) {
	f, ok := new(big.Float).SetString(fil)
	if !ok {
		return "", fmt.Errorf("invalid FIL amount: %s", fil)
	}

	attoPerFil := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	attoFloat := new(big.Float).Mul(f, attoPerFil)

	attoInt, _ := attoFloat.Int(nil)
	return attoInt.String(), nil
}

// FormatFIL converts an attoFIL string to a human-readable FIL string with 6 decimals.
func FormatFIL(attoFil string) (string, error) {
	attoInt := new(big.Int)
	if _, ok := attoInt.SetString(attoFil, 10); !ok {
		return "", fmt.Errorf("invalid attoFIL amount: %s", attoFil)
	}

	attoPerFil := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	filFloat := new(big.Float).Quo(
		new(big.Float).SetInt(attoInt),
		new(big.Float).SetInt(attoPerFil),
	)

	return filFloat.Text('f', 6), nil
}
