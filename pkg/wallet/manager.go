package wallet

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/filament-labs/filament/pkg/database"
	"github.com/filament-labs/filament/pkg/util"
	"github.com/filecoin-project/go-address"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
)

type Manager struct {
	store     database.Store
	wallets   map[string]*Wallet
	masterKey []byte
	locked    bool
}

const walletsDBPrefix = "wallets_"

func NewManager(store database.Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) Lock() {
	m.wallets = make(map[string]*Wallet)
	m.masterKey = nil
	m.locked = true
}

func (m *Manager) Unlock(password string) error {
	m.masterKey = pbkdf2.Key([]byte(password), []byte("filament-global-salt"), 4096, 36, sha256.New)
	m.locked = false

	_, err := m.LoadWallets()
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) CreateWallet(name string) (*Wallet, error) {
	if m.locked {
		return nil, fmt.Errorf("wallet manager is locked")
	}

	// Generate BIP-39 mnemonic
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to create mnemonic: %w", err)
	}

	// Derive seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Derive Filecoin secp256k1 key: m/44'/461'/0'/0/0
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	path := []uint32{44, 461, 0, 0, 0}
	privateKey := masterKey
	for _, n := range path {
		privateKey, err = privateKey.NewChildKey(n)
		if err != nil {
			return nil, fmt.Errorf("failed to derive child key: %w", err)
		}
	}

	pubKey := privateKey.PublicKey().Key
	addr, err := address.NewSecp256k1Address(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	// Generate random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// ✅ Set IsDefault before saving
	isDefault := len(m.wallets) == 0

	wallet := &Wallet{
		ID:        uuid.New().String(),
		Name:      name,
		Address:   addr.String(),
		Balance:   &Balance{}, // avoid nil
		IsDefault: isDefault,
		CreatedAt: time.Now().Unix(),
		salt:      salt,
	}

	// Derive encryption key from masterKey and wallet-specific salt
	encryptionKey := pbkdf2.Key(m.masterKey, wallet.salt, 4096, 32, sha256.New)

	// Encrypt mnemonic
	if err := util.EncryptData([]byte(mnemonic), encryptionKey, &wallet.encryptedSeed, &wallet.nonceSeed); err != nil {
		return nil, fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}

	// Encrypt private key
	if err := util.EncryptData(privateKey.Key, encryptionKey, &wallet.encryptedKey, &wallet.nonceKey); err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Save wallet
	key := walletsDBPrefix + util.HyphenateAndLower(wallet.Name)
	data, err := util.Encode(wallet)
	if err != nil {
		return nil, fmt.Errorf("error encoding wallet: %w", err)
	}

	if err := m.store.Save(key, data); err != nil {
		return nil, fmt.Errorf("error saving wallet to database: %w", err)
	}

	m.wallets[wallet.ID] = wallet
	return wallet, nil
}

func (m *Manager) LoadWallets() ([]*Wallet, error) {
	if m.locked {
		return nil, errors.New("wallet manager is locked")
	}

	var wallets []*Wallet
	_, err := m.store.GetMany(walletsDBPrefix, func(val []byte) (any, error) {
		var w Wallet
		if err := util.Decode(val, &w); err != nil {
			return nil, err
		}
		wallets = append(wallets, &w)
		return w, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching wallets: %w", err)
	}

	m.wallets = make(map[string]*Wallet, len(wallets))
	for _, wal := range wallets {
		// Validate decryption
		verified := false
		for _, wal := range wallets {
			if !verified {
				if _, err := wal.GetPrivateKey(m.masterKey); err != nil {
					return nil, fmt.Errorf("invalid password: failed to decrypt wallet %s: %w", wal.ID, err)
				}
				verified = true
			}
			m.wallets[wal.ID] = wal
		}
		m.wallets[wal.ID] = wal
	}

	return wallets, nil
}

// RecoverWallet recreates a wallet from a given BIP-39 mnemonic.
// It encrypts the key material using the user's master password and saves it to the database.
func (m *Manager) RecoverWallet(name, mnemonic string) (*Wallet, error) {
	if m.locked {
		return nil, fmt.Errorf("wallet manager is locked")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic phrase")
	}

	// Derive seed from mnemonic (no passphrase)
	seed := bip39.NewSeed(mnemonic, "")

	// Derive Filecoin secp256k1 key: m/44'/461'/0'/0/0
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %w", err)
	}

	path := []uint32{44, 461, 0, 0, 0}
	privateKey := masterKey
	for _, n := range path {
		privateKey, err = privateKey.NewChildKey(n)
		if err != nil {
			return nil, fmt.Errorf("failed to derive child key: %w", err)
		}
	}

	pubKey := privateKey.PublicKey().Key
	addr, err := address.NewSecp256k1Address(pubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	// Generate new salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	isDefault := len(m.wallets) == 0

	wallet := &Wallet{
		ID:        uuid.New().String(),
		Name:      name,
		Address:   addr.String(),
		Balance:   &Balance{},
		IsDefault: isDefault,
		CreatedAt: time.Now().Unix(),
		salt:      salt,
	}

	// Derive encryption key from masterKey and wallet salt
	encryptionKey := pbkdf2.Key(m.masterKey, wallet.salt, 4096, 32, sha256.New)

	// Encrypt mnemonic and private key again
	if err := util.EncryptData([]byte(mnemonic), encryptionKey, &wallet.encryptedSeed, &wallet.nonceSeed); err != nil {
		return nil, fmt.Errorf("failed to encrypt mnemonic: %w", err)
	}
	if err := util.EncryptData(privateKey.Key, encryptionKey, &wallet.encryptedKey, &wallet.nonceKey); err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Save to DB
	key := walletsDBPrefix + util.HyphenateAndLower(wallet.Name)
	data, err := util.Encode(wallet)
	if err != nil {
		return nil, fmt.Errorf("error encoding wallet: %w", err)
	}
	if err := m.store.Save(key, data); err != nil {
		return nil, fmt.Errorf("error saving recovered wallet: %w", err)
	}

	m.wallets[wallet.ID] = wallet
	return wallet, nil
}
