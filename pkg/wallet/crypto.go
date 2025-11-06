package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-address"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
)

const (
	// Standard BIP39 derivation path for Filecoin
	defaultDerivationPath = "m/44'/461'/0'/0/0"
	saltSize              = 32
	keySize               = 32
	nonceSize             = 12
	pbkdf2Iterations      = 100000
)

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

// DeriveAddressFromPrivateKey derives Filecoin addresses from private key
func DeriveAddressFromPrivateKey(privKey *ecdsa.PrivateKey) (map[string]string, error) {
	pubBytes := crypto.FromECDSAPub(&privKey.PublicKey)

	// Create f1 (secp256k1) address
	f1Addr, err := address.NewSecp256k1Address(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("create f1 address: %w", err)
	}

	return map[string]string{
		"f1":  f1Addr.String(),
		"hex": hex.EncodeToString(pubBytes),
	}, nil
}

// EncryptMnemonic encrypts a mnemonic phrase with a passphrase
func EncryptMnemonic(mnemonic, passphrase string) ([]byte, error) {
	if passphrase == "" {
		return nil, ErrInvalidPassphrase
	}

	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	// Derive key using PBKDF2
	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, keySize, sha256.New)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	// Encrypt mnemonic
	ciphertext := gcm.Seal(nil, nonce, []byte(mnemonic), nil)

	// Format: salt + nonce + ciphertext
	result := make([]byte, 0, saltSize+nonceSize+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptMnemonic decrypts an encrypted mnemonic with a passphrase
func DecryptMnemonic(encrypted []byte, passphrase string) (string, error) {
	if passphrase == "" {
		return "", ErrInvalidPassphrase
	}

	if len(encrypted) < saltSize+nonceSize {
		return "", fmt.Errorf("encrypted data too short")
	}

	// Extract salt, nonce, and ciphertext
	salt := encrypted[:saltSize]
	nonce := encrypted[saltSize : saltSize+nonceSize]
	ciphertext := encrypted[saltSize+nonceSize:]

	// Derive key using PBKDF2
	key := pbkdf2.Key([]byte(passphrase), salt, pbkdf2Iterations, keySize, sha256.New)

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// GenerateMnemonic generates a new BIP39 mnemonic phrase
func GenerateMnemonic(bits int) (string, error) {
	if bits != 128 && bits != 256 {
		bits = 128 // default to 12 words
	}

	entropy, err := bip39.NewEntropy(bits)
	if err != nil {
		return "", fmt.Errorf("generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// ValidateMnemonic checks if a mnemonic phrase is valid
func ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// MnemonicToSeed converts a mnemonic to a seed
func MnemonicToSeed(mnemonic, password string) []byte {
	return bip39.NewSeed(mnemonic, password)
}
