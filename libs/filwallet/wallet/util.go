package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/argon2"
)

// deriveMasterKey turns a human password into a high-entropy 32-byte hex string.
func deriveMasterKey(password string, salt []byte) []byte {
	// Argon2id used to protect against GPU brute-forcing.
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
}

func encryptAESGCM(plaintext, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key[:32]) // Use first 32 bytes
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decryptAESGCM(ciphertext, key []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key[:32])
	gcm, _ := cipher.NewGCM(block)
	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("invalid ciphertext")
	}
	nonce, ct := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

// derivePrivateKeyFromSeed derives an ECDSA private key from a BIP39 seed
func derivePrivateKeyFromSeed(seed []byte) (*ecdsa.PrivateKey, error) {
	// Use first 32 bytes of seed as private key material
	if len(seed) < 32 {
		return nil, fmt.Errorf("seed too short")
	}

	// Hash the seed to get deterministic private key
	hash := sha256.Sum256(seed[:32])
	privKey, err := crypto.ToECDSA(hash[:])
	if err != nil {
		return nil, fmt.Errorf("convert to ecdsa: %w", err)
	}

	return privKey, nil
}

func wipeECDSA(priv *ecdsa.PrivateKey) {
	if priv == nil || priv.D == nil {
		return
	}
	dBytes := priv.D.Bits()
	for i := range dBytes {
		dBytes[i] = 0
	}
}

func importKeyJSON(masterKey []byte, privKey *ecdsa.PrivateKey, dataDir string) ([]byte, error) {
	ks := keystore.NewKeyStore(dataDir, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privKey, string(masterKey))
	if err != nil {
		if errors.Is(err, keystore.ErrAccountAlreadyExists) {
			return nil, ErrWalletAlreadyExists
		}
		return nil, fmt.Errorf("import ecdsa: %w", err)
	}

	keyJSON, err := ks.Export(account, string(masterKey), string(masterKey))
	if err != nil {
		return nil, fmt.Errorf("export keystore: %w", err)
	}

	return keyJSON, nil
}
