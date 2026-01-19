package filwallet

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/awnumar/memguard"
	"github.com/codemaestro64/filament/libs/filwallet/wallet"
)

type sessionState struct {
	vault     map[int]*memguard.Enclave
	expiresAt time.Time
}

type Manager struct {
	cfg       *Config
	store     Store
	rpcClient *RPCClient
	session   *sessionState
	mu        sync.RWMutex
}

func NewManager(ctx context.Context, store Store, cfg *Config) (*Manager, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("initialize wallet manager: %w", err)
	}

	rpcClient, err := NewRPCClient(ctx, cfg.RPCEndpoint, cfg.RPCToken)
	if err != nil {
		return nil, fmt.Errorf("initialize rpc client: %w", err)
	}

	m := &Manager{
		cfg:       cfg,
		rpcClient: rpcClient,
		store:     store,
		session: &sessionState{
			vault:     make(map[int]*memguard.Enclave),
			expiresAt: time.Now().Add(30 * time.Minute),
		},
	}

	m.startSessionJanitor(ctx)

	return m, nil
}

func (m *Manager) importWallet(ctx context.Context, mnemonic, walletName, password string) (*wallet.Wallet, error) {
	newWallet, err := wallet.CreateNew(m.cfg.DataDir, mnemonic, walletName, password)
	if err != nil {
		return nil, fmt.Errorf("create wallet: %w", err)
	}

	dbWallet, err := m.store.SaveWallet(ctx, SaveWalletParams{
		KeyJSON:       newWallet.EncryptedKeyJSON,
		EncryptedSeed: newWallet.EncryptedMnemonic,
		Addresses:     newWallet.Addresses,
		Name:          newWallet.Name,
		Salt:          newWallet.Salt,
	})
	if err != nil {
		return nil, fmt.Errorf("save wallet: %w", err)
	}

	newWallet.ID = dbWallet.ID
	err = m.UnlockWallet(ctx, newWallet.ID, password)
	if err != nil {
		return nil, fmt.Errorf("unlock wallet: %w", err)
	}

	return newWallet, nil
}

func (m *Manager) UnlockWallet(ctx context.Context, walletID int, password string) error {
	m.mu.RLock()
	if _, exists := m.session.vault[walletID]; exists {
		m.mu.RUnlock()
		return nil // Already unlocked, no work needed
	}
	m.mu.RUnlock()

	wallet, err := m.store.FindWallet(ctx, walletID)
	if err != nil {
		return fmt.Errorf("find wallet: %w", err)
	}

	enclave, err := wallet.Unlock(password)
	if err != nil {
		return fmt.Errorf("unlock wallet: %w", err)
	}

	m.mu.Lock()
	m.session.vault[wallet.ID] = enclave
	m.mu.Unlock()

	expireDuration := time.Minute * time.Duration(m.cfg.SessionTimeout)
	m.session.expiresAt = time.Now().Add(expireDuration)

	return nil
}

func (m *Manager) UnlockAllWallets(ctx context.Context, password string) error {
	wallets, err := m.store.GetWallets(ctx)
	if err != nil {
		return fmt.Errorf("get wallets: %w", err)
	}

	tempVault := make(map[int]*memguard.Enclave)
	for _, wallet := range wallets {
		enclave, err := wallet.Unlock(password)
		if err != nil {
			return fmt.Errorf("unlock wallet %d: %w", wallet.ID, err)
		}
		tempVault[wallet.ID] = enclave
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for id, enclave := range tempVault {
		m.session.vault[id] = enclave
	}

	expireDuration := time.Minute * time.Duration(m.cfg.SessionTimeout)
	m.session.expiresAt = time.Now().Add(expireDuration)

	return nil
}

func (m *Manager) WalletsCount(ctx context.Context) (int, error) {
	count, err := m.store.CountWallets(ctx)
	if err != nil {
		return 0, fmt.Errorf("count db wallets: %w", err)
	}

	return count, nil
}

func (m *Manager) RecoverWallet(ctx context.Context, seedWords, walletName, password string) (*wallet.Wallet, error) {
	if !ValidateMnemonic(seedWords) {
		return nil, ErrInvalidSeedPhrase
	}

	if password == "" {
		return nil, ErrInvalidPassword
	}

	if walletName == "" {
		return nil, ErrInvalidWalletName
	}

	return m.importWallet(ctx, seedWords, walletName, password)
}

func (m *Manager) CreateWallet(ctx context.Context, walletName, password string) (*wallet.Wallet, string, error) {
	if password == "" {
		return nil, "", ErrInvalidPassword
	}

	if walletName == "" {
		return nil, "", ErrInvalidWalletName
	}

	mnemonic, err := GenerateMnemonic(128)
	if err != nil {
		return nil, "", fmt.Errorf("generate seed words: %w", err)
	}

	wallet, err := m.importWallet(ctx, mnemonic, walletName, password)
	if err != nil {
		return nil, "", fmt.Errorf("create wallet: %w", err)
	}

	return wallet, mnemonic, nil
}

func (m *Manager) LockWallets() {
	m.mu.Lock()
	defer m.mu.RUnlock()

	m.session.vault = make(map[int]*memguard.Enclave)
	m.session.expiresAt = time.Time{}
}

func (m *Manager) startSessionJanitor(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				m.LockWallets()
				return
			case <-ticker.C:
				m.mu.Lock()
				if !m.session.expiresAt.IsZero() && time.Now().After(m.session.expiresAt) {
					m.mu.Unlock()
					m.LockWallets()
				} else {
					m.mu.Unlock()
				}
			}
		}
	}()
}
