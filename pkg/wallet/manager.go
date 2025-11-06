package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
)

// Manager coordinates wallet lifecycle, persistence, and unlocking.
// It holds in-memory references to loaded wallets and delegates storage to a *store.
type Manager struct {
	wallets   map[string]Wallet // walletID → Wallet
	store     *store
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
func NewManager(db *badger.DB, opts ...ManagerOption) (*Manager, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}

	store := newStore(db)

	loaded, err := store.listWallets()
	if err != nil {
		return nil, fmt.Errorf("failed to load wallets from DB: %w", err)
	}

	wallets := make(map[string]Wallet, len(loaded))
	for _, w := range loaded {
		wallets[w.ID] = w
	}

	m := &Manager{
		store:   store,
		wallets: wallets,
	}

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

// isWalletUnlocked reports whether the wallet with the given ID is currently unlocked in memory.
func (m *Manager) isWalletUnlocked(walletID string) bool {
	w, exists := m.wallets[walletID]
	return exists && w.key != nil
}

func (m *Manager) HasWallets() bool {
	return len(m.wallets) > 0
}

// UnlockWallets attempts to decrypt and unlock every locked wallet using the supplied keystore passphrase.
// Already-unlocked wallets are skipped. If any wallet fails to decrypt, the whole operation aborts.
func (m *Manager) UnlockWallets(keystorePassphrase string) error {
	for id, w := range m.wallets {
		if m.isWalletUnlocked(id) {
			continue
		}

		key, err := keystore.DecryptKey(w.KeyJSON, keystorePassphrase)
		if err != nil {
			return fmt.Errorf("failed to unlock wallet %s: %w", id, err)
		}

		// Mutate the copy in the map
		w.key = key
		m.wallets[id] = w
	}
	return nil
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

	// Save wallet
	if err := m.store.saveWallet(*wallet); err != nil {
		return nil, err
	}

	m.wallets[wallet.ID] = *wallet

	return wallet, nil
}
