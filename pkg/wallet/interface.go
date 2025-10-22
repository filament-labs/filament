package wallet

import (
	"github.com/filament-labs/filament/internal/models"
)

type WalletManager interface {
	LoadWallets() ([]models.WalletInfo, error)
	CreateWallet(name, password string) (*models.WalletInfo, error)
	GetWallet(walletName string) (*Wallet, error)
	//sUnlock(password string) error
	//GetPrivateKey(walletName string) (*ecdsa.PrivateKey, error)
	//SignTransaction(walletName string, tx *models.Transaction) (*models.SignedTransaction, error)
	//SetDefaultWallet(walletName string) error
}
