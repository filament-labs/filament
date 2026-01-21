package wallet

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/awnumar/memguard"
	"github.com/codemaestro64/filament/libs/filwallet/address"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	ID                int
	IsDefault         bool
	Name              string
	Addresses         []address.Address
	Salt              []byte
	EncryptedKeyJSON  []byte
	EncryptedMnemonic []byte
	CreatedAt         time.Time
}

func CreateNew(dataDir string, mnemonic, walletName, password string) (*Wallet, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	seed := bip39.NewSeed(mnemonic, "")
	defer memguard.WipeBytes(seed)

	privKey, err := derivePrivateKeyFromSeed(seed)
	if err != nil {
		return nil, fmt.Errorf("derive private key from seed: %w", err)
	}
	defer wipeECDSA(privKey)

	masterKey := deriveMasterKey(password, salt)
	defer memguard.WipeBytes(masterKey)

	keyJSON, err := importKeyJSON(masterKey, privKey, dataDir)
	if err != nil {
		return nil, fmt.Errorf("import keyJSON: %w", err)
	}

	encryptedMnemonic, err := encryptAESGCM([]byte(mnemonic), masterKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt mnemonic: %w", err)
	}

	addresses, err := address.DeriveAddressesFromPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("derive addresses: %w", err)
	}

	return &Wallet{
		Name:              walletName,
		Salt:              salt,
		Addresses:         addresses,
		EncryptedKeyJSON:  keyJSON,
		EncryptedMnemonic: encryptedMnemonic,
	}, nil
}

// Unlock handles the decryption logic internal to a wallet's data.
func (w *Wallet) Unlock(password string) (*memguard.Enclave, error) {
	masterKey := deriveMasterKey(password, w.Salt)
	defer memguard.WipeBytes(masterKey)

	key, err := keystore.DecryptKey(w.EncryptedKeyJSON, string(masterKey))
	if err != nil {
		return nil, fmt.Errorf("decrypt wallet key: %w", err)
	}

	privBytes := crypto.FromECDSA(key.PrivateKey)
	enclave := memguard.NewEnclave(privBytes)

	// Clean up
	wipeECDSA(key.PrivateKey)
	memguard.WipeBytes(privBytes)

	return enclave, nil
}

func (w *Wallet) DecryptSeedPhrase(password string) (string, error) {
	if len(w.EncryptedMnemonic) == 0 {
		return "", errors.New("no seed phrase stored for this wallet")
	}

	masterKey := deriveMasterKey(password, w.Salt)
	defer memguard.WipeBytes(masterKey)

	mnemonicBytes, err := decryptAESGCM(w.EncryptedMnemonic, masterKey)
	if err != nil {
		return "", fmt.Errorf("decrypt seed phrase: %w", err)
	}

	buf := memguard.NewBufferFromBytes(mnemonicBytes)
	//defer buf.Destroy()

	return buf.String(), nil
}
