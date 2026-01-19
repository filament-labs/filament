package filwallet

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/awnumar/memguard"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
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

func (m *Manager) WalletsLocked() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check if session exists at all
	if m.session == nil || m.session.vault == nil {
		return true
	}

	// Check if the session has expired
	if !m.session.expiresAt.IsZero() && time.Now().After(m.session.expiresAt) {
		return true
	}

	// Locked if the vault is empty
	return len(m.session.vault) == 0
}

func (m *Manager) CountWallets(ctx context.Context) (int, error) {
	count, err := m.store.CountWallets(ctx)
	if err != nil {
		return 0, fmt.Errorf("count db wallets: %w", err)
	}

	return count, nil
}

func (m *Manager) unlockWallet(wallet *Wallet, password string) error {
	// Derive master key that serves as real password for all storage
	// User provided password serves as entropy for the derived master key
	masterKey := deriveMasterKey(password, wallet.Salt)

	// Decrypt kay
	key, err := keystore.DecryptKey(wallet.EncryptedKeyJSON, masterKey)
	if err != nil {
		return fmt.Errorf("decrypt wallet %d: %w", wallet.ID, err)
	}

	// Move to Memguard encalve immediately
	privBytes := crypto.FromECDSA(key.PrivateKey)
	m.session.vault[wallet.ID] = memguard.NewEnclave(privBytes)

	// Wipe sensitive data from memory (heap)
	wipeECDSA(key.PrivateKey)
	memguard.WipeBytes(privBytes)

	return nil
}

func (m *Manager) UnlockWallet(ctx context.Context, walletID int, password string) error {
	wallet, err := m.store.FindWallet(ctx, walletID)
	if err != nil {
		return fmt.Errorf("unlock wallet: %w", err)
	}

	return m.unlockWallet(wallet, password)
}

func (m *Manager) UnlockWallets(ctx context.Context, password string) error {
	wallets, err := m.store.GetWallets(ctx)
	if err != nil {
		return fmt.Errorf("fetch wallets: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, wallet := range wallets {
		if err := m.unlockWallet(wallet, password); err != nil {
			return err
		}
	}
	m.session.expiresAt = time.Now().Add(m.cfg.SessionDuration)

	return nil
}

func (m *Manager) GetKey(walletID int) (*memguard.LockedBuffer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || time.Now().After(m.session.expiresAt) {
		return nil, ErrSessionExpired
	}

	enclave, ok := m.session.vault[walletID]
	if !ok {
		return nil, ErrNotFound
	}

	// privKey, err := crypto.ToECDSA(buffer.Bytes())

	return enclave.Open()
}

func (m *Manager) RevealSeedPhrase(ctx context.Context, walletID int, password string) (string, error) {
	wallet, err := m.store.FindWallet(ctx, walletID)
	if err != nil {
		return "", fmt.Errorf("reveal seed phrase: %w", err)
	}

	if len(wallet.EncryptedSeed) == 0 {
		return "", errors.New("no seed phrase stored for this wallet")
	}

	// Derive master key
	masterKey := deriveMasterKey(password, wallet.Salt)
	defer memguard.WipeBytes([]byte(masterKey))

	seedBytes, err := decryptAESGCM(wallet.EncryptedSeed, []byte(masterKey))
	if err != nil {
		return "", fmt.Errorf("decrypt seed phrase: %w", err)
	}

	seedBuf := memguard.NewBufferFromBytes(seedBytes)
	defer seedBuf.Destroy()
	memguard.WipeBytes(seedBytes)

	return seedBuf.String(), nil
}

func (m *Manager) importWallet(ctx context.Context, seedWords, walletName, password string) (*Wallet, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	seed := bip39.NewSeed(seedWords, "")
	privKey, err := derivePrivateKeyFromSeed(seed)
	if err != nil {
		return nil, fmt.Errorf("derive private key: %w", err)
	}
	defer wipeECDSA(privKey)
	defer memguard.WipeBytes(seed)

	masterKey := deriveMasterKey(password, salt)
	defer memguard.WipeBytes([]byte(masterKey))

	ks := keystore.NewKeyStore(m.cfg.DataDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privKey, masterKey)
	if err != nil {
		if errors.Is(err, keystore.ErrAccountAlreadyExists) {
			return nil, ErrWalletExists
		}
		return nil, fmt.Errorf("import wallet: %w", err)
	}

	keyJSON, err := ks.Export(account, masterKey, masterKey)
	if err != nil {
		return nil, fmt.Errorf("export keystore: %w", err)
	}

	addresses, err := DeriveAddressesFromPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("derive addresses: %w", err)
	}

	// Encrypt the mnemonic (seedWords)
	encryptedSeedPhrase, err := encryptAESGCM([]byte(seedWords), []byte(masterKey))
	if err != nil {
		return nil, err
	}

	dbWallet, err := m.store.SaveWallet(ctx, SaveWalletParams{
		KeyJSON:       keyJSON,
		EncryptedSeed: encryptedSeedPhrase,
		Addresses:     addresses,
		Name:          walletName,
		Salt:          salt,
		Password:      password,
	})
	if err != nil {
		return nil, fmt.Errorf("save wallet: %w", err)
	}

	err = m.UnlockWallet(ctx, dbWallet.ID, password)
	if err != nil {
		return nil, fmt.Errorf("unlock wallet: %w", err)
	}

	return &Wallet{
		ID:         dbWallet.ID,
		IsDefault:  dbWallet.IsDefault,
		SeedPhrase: seedWords,
		Addresses:  addresses,
		CreatedAt:  dbWallet.CreatedAt,
	}, nil
}

func (m *Manager) CreateWallet(ctx context.Context, walletName, password string) (*Wallet, error) {
	if password == "" {
		return nil, ErrInvalidPassword
	}

	if walletName == "" {
		return nil, ErrInvalidWalletName
	}

	seedWords, err := GenerateMnemonic(128)
	if err != nil {
		return nil, fmt.Errorf("generate seed words: %w", err)
	}

	return m.importWallet(ctx, seedWords, walletName, password)
}

func (m *Manager) RecoverWallet(ctx context.Context, seedWords, walletName, password string) (*Wallet, error) {
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

func (m *Manager) LockWallets() {
	m.mu.Lock()
	defer m.mu.RUnlock() // Ensure we use RUnlock if we were just reading, or Unlock if writing

	m.session.vault = make(map[int]*memguard.Enclave)
	m.session.expiresAt = time.Time{}
}

// Close gracefully shuts down the RPC and wipes sessions.
func (m *Manager) Close() {
	m.rpcClient.Close()
}
