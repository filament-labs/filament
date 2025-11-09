package dto

import "time"

type WalletRequest struct {
	WalletID string `json:"wallet_id"`
}

type WalletResponse struct {
	ID        string             `json:"id"`
	IsDefault bool               `json:"is_default"`
	Name      string             `json:"name"`
	Addresses map[string]string  `json:"addresses"`
	Balance   GetBalanceResponse `json:"balance"`
	CreatedAt time.Time          `json:"created_at"`
}

type WalletsRequest struct{}
type WalletsResponse struct {
	Locked  bool             `json:"locked"`
	IsInit  bool             `json:"is_init"`
	Wallets []WalletResponse `json:"wallets"`
}

type UnlockWalletsRequest struct {
	Password string `json:"password"`
}

type UnlockWalletsResponse struct{}
