package wallet

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filament-labs/filament/pkg/util"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	ID        string   `json:"id"`         // Unique wallet identifier
	IsDefault bool     `json:"is_default"` // Is user's default wallet
	Name      string   `json:"name"`       // User-provided wallet name
	Address   string   `json:"address"`    // Filecoin f1 address
	Balance   *Balance `json:"balance"`    // Wallet balance
	CreatedAt int64    `json:"created_at"` // Timestamp

	encryptedSeed []byte // Encrypted BIP-39 mnemonic
	encryptedKey  []byte // Encrypted secp256k1 private key
	nonceSeed     []byte // Nonce for seed encryption
	nonceKey      []byte // Nonce for key encryption
	salt          []byte // Unique salt for this wallet
}

func (w *Wallet) GetPrivateKey(masterKey []byte) (*ecdsa.PrivateKey, error) {
	privBytes, err := util.DecryptData(w.encryptedKey, w.nonceKey, masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %w", err)
	}

	privKey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privKey, nil
}

func (w *Wallet) GetMnemonic(masterKey []byte) (string, error) {
	seedBytes, err := util.DecryptData(w.encryptedSeed, w.nonceSeed, masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt mnemonic: %w", err)
	}

	mnemonic := string(seedBytes)
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic")
	}

	return mnemonic, nil
}
