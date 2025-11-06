package dto

type WalletCountRequest struct{}
type WalletCountResponse struct {
	WalletCount int `json:"wallet_count"`
}

type WalletRequest struct {
	WalletID   string `json:"wallet_id"`
	WalletName string `json:"wallet_name"`
}

type WalletResponse struct {
	ID        string             `json:"id"`
	IsDefault bool               `json:"is_default"`
	Name      string             `json:"name"`
	Addresses map[string]string  `json:"addresses"`
	Balance   GetBalanceResponse `json:"balance"`
	CreatedAt int64              `json:"created_at"`
}

type WalletsRequest struct{}
type WalletsResponse struct {
	Locked  bool             `json:"locked"`
	Wallets []WalletResponse `json:"wallets"`
}

type UnlockWalletsRequest struct {
	Password string `json:"password"`
}

type UnlockWalletsResponse struct{}
