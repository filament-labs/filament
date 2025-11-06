package wallet

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
)

type Manager struct {
	rpcClient *RPCClient
}

type ManagerOption func(*Manager) error

// WithRPC allows setting a custom RPC endpoint and optional token
func WithOptions(url, token string) ManagerOption {
	return func(m *Manager) error {
		if url == "" {
			url = "https://filfox.info/rpc/v1"
		}
		rpcClient, err := NewRPCClient(RPCConfig{
			Endpoint: url,
			Token:    token,
		})
		if err != nil {
			return fmt.Errorf("failed to initialize RPC client: %w", err)
		}
		m.rpcClient = rpcClient
		return nil
	}
}

// NewManager initializes a Manager with Badger DB and optional configurations
func NewManager(opts ...ManagerOption) (*Manager, error) {
	m := &Manager{}

	// Apply options
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}

	// If no RPC option was applied, initialize with default
	if m.rpcClient == nil {
		rpcClient, err := NewRPCClient(RPCConfig{
			Endpoint: "https://filfox.info/rpc/v1",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize default RPC client: %w", err)
		}
		m.rpcClient = rpcClient
	}

	return m, nil
}

func (m *Manager) UnlockWallet(wallet *Wallet, keystorePassphrase string) (*keystore.Key, error) {
	key, err := keystore.DecryptKey(wallet.KeyJSON, keystorePassphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to unlock wallet %s: %w", wallet.ID, err)
	}

	return key, nil
}

// CreateWallet generates a brand-new wallet, persists it, and returns the in-memory instance (locked).
func (m *Manager) CreateWallet(ctx context.Context, walletName, keystorePassphrase string) (*Wallet, error) {
	if keystorePassphrase == "" {
		return nil, ErrInvalidPassphrase
	}

	if walletName == "" {
		return nil, ErrInvalidWalletName
	}

	// Generate new mnemonic (12 words)
	mnemonic, err := GenerateMnemonic(128)
	if err != nil {
		return nil, fmt.Errorf("error generating mnemonic: %w", err)
	}

	return m.createWalletFromMnemonic(mnemonic, walletName, keystorePassphrase, false)
}

// RecoverWallet recovers a wallet from a mnemonic phrase
func (m *Manager) RecoverWallet(ctx context.Context, mnemonic, name, keystorePassphrase string) (*Wallet, error) {
	if !ValidateMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	if keystorePassphrase == "" {
		return nil, ErrInvalidPassphrase
	}

	return m.createWalletFromMnemonic(mnemonic, name, keystorePassphrase, true)
}

// createWalletFromMnemonic creates a wallet from a mnemonic phrase
func (m *Manager) createWalletFromMnemonic(mnemonic, name, keystorePassphrase string, recovered bool) (*Wallet, error) {
	// Generate seed from mnemonic
	seed := MnemonicToSeed(mnemonic, "")

	// Derive private key
	privKey, err := derivePrivateKeyFromSeed(seed)
	if err != nil {
		return nil, fmt.Errorf("derive private key: %w", err)
	}

	// Create keystore encryption
	ks := keystore.NewKeyStore("", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privKey, keystorePassphrase)
	if err != nil {
		return nil, fmt.Errorf("import ecdsa: %w", err)
	}

	keyJSON, err := ks.Export(account, keystorePassphrase, keystorePassphrase)
	if err != nil {
		return nil, fmt.Errorf("export keystore: %w", err)
	}

	// Derive addresses
	addrs, err := DeriveAddressFromPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("derive addresses: %w", err)
	}

	now := time.Now()
	wallet := &Wallet{
		ID:        uuid.NewString(),
		Name:      name,
		Mnemonic:  mnemonic,
		KeyJSON:   keyJSON,
		Addrs:     addrs,
		Meta:      make(map[string]string),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if recovered {
		wallet.Meta["recovered"] = "true"
	}

	return wallet, nil
}
