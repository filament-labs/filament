package models

type WalletInfo struct {
	ID        string
	Name      string
	Address   string
	IsDefault bool
	Balance   *Balance
}
