package filwallet

import (
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

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
