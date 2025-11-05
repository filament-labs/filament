package wallet

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/google/uuid"
)

type WalletManager struct {
	wallets map[string]*Wallet
	store   WalletStore
}

func NewManager(walletStore WalletStore) (*WalletManager, error) {
	wm := &WalletManager{
		store: walletStore,
	}

	return wm, nil
}

func (wm *WalletManager) CreateWallet(ctx context.Context, name, passphrase string) (*Wallet, error) {
	tsDir := ""
	ks := keystore.NewKeyStore(tsDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acct, err := ks.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}

	keyJSON, err := ks.Export(acct, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	decryptedKey, err := keystore.DecryptKey(keyJSON, passphrase)
	if err != nil {
		return nil, err
	}

	derivedAddrs, err := DeriveFromECDSA(decryptedKey.PrivateKey)
	if err != nil {
		return nil, err
	}

	walletRec := &Wallet{
		ID:      uuid.NewString(),
		Name:    name,
		KeyJSON: keyJSON,
		Addrs: map[string]string{
			"f1":  derivedAddrs.F1.String(),
			"f4":  derivedAddrs.F4.String(),
			"hex": derivedAddrs.Hex,
		},
		Meta: map[string]string{
			"name":       name,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}

	// TODO save wallet
	wm.wallets[name] = walletRec
	return walletRec, nil
}

func (wm *WalletManager) LoadWallets() error {
	return nil
}

func (wm *WalletManager) GetWallet(walletID string) (*Wallet, error) {
	if wallet, ok := wm.wallets[walletID]; ok {
		return wallet, nil
	}

	return nil, ErrWalletNotFound
}
