package model

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filament-labs/filament/internal/dto"
	"github.com/filament-labs/filament/pkg/util"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	dto.GetWalletResponse
	EncryptedSeed []byte // Encrypted BIP-39 mnemonic
	EncryptedKey  []byte // Encrypted secp256k1 private key
	NonceSeed     []byte // Nonce for seed encryption
	NonceKey      []byte // Nonce for key encryption
	Salt          []byte // Unique salt for this wallet
}

func (w *Wallet) GetPrivateKey(masterKey []byte) (*ecdsa.PrivateKey, error) {
	privBytes, err := util.DecryptData(w.EncryptedKey, w.NonceKey, masterKey)
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
	seedBytes, err := util.DecryptData(w.EncryptedSeed, w.NonceSeed, masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt mnemonic: %w", err)
	}

	mnemonic := string(seedBytes)
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", fmt.Errorf("invalid mnemonic")
	}

	return mnemonic, nil
}
