package filwallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	"github.com/awnumar/memguard"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-address"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/argon2"
)

const (
	evmNamespace uint64 = 10
)

func wipeECDSA(priv *ecdsa.PrivateKey) {
	if priv == nil || priv.D == nil {
		return
	}
	dBytes := priv.D.Bits()
	for i := range dBytes {
		dBytes[i] = 0
	}
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

// deriveMasterKey turns a human password into a high-entropy 32-byte hex string.
func deriveMasterKey(password string, salt []byte) string {
	// Argon2id is used to protect against GPU brute-forcing.
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	defer memguard.WipeBytes(hash)

	return hex.EncodeToString(hash)
}

// DeriveAddressesFromPrivateKey returns all standard Filecoin address formats
func DeriveAddressesFromPrivateKey(privKey *ecdsa.PrivateKey) ([]*pb.Address, error) {
	// Get raw private key bytes (32 bytes)
	//privBytes := crypto.FromECDSA(privKey)

	// f1 address (legacy secp256k1)
	pubBytes := crypto.CompressPubkey(&privKey.PublicKey)
	f1Addr, err := address.NewSecp256k1Address(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("create f1 address: %w", err)
	}

	// f4 address (delegated/eth address)
	ethAddr := crypto.PubkeyToAddress(privKey.PublicKey) // Ethereum address
	f4Addr, err := address.NewDelegatedAddress(evmNamespace, ethAddr.Bytes())
	if err != nil {
		return nil, fmt.Errorf("create f4 address: %w", err)
	}

	return []*pb.Address{
		{Type: pb.AddressType_ADDRESS_TYPE_F1, Value: f1Addr.String()},
		{Type: pb.AddressType_ADDRESS_TYPE_F4, Value: f4Addr.String()},
		{Type: pb.AddressType_ADDRESS_TYPE_0X, Value: ethAddr.Hex()},
	}, nil
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
