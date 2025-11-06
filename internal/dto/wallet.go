package dto

type GetWalletRequest struct {
	WalletID   string `json:"wallet_id"`
	WalletName string `json:"wallet_name"`
}

type GetWalletResponse struct {
	ID        string             `json:"id"`
	IsDefault bool               `json:"is_default"`
	Name      string             `json:"name"`
	Addresses map[string]string  `json:"addresses"`
	Balance   GetBalanceResponse `json:"balance"`
	CreatedAt int64              `json:"created_at"`
}

type GetWalletsRequest struct{}
type GetWalletsResponse struct {
	Locked  bool                `json:"locked"`
	Wallets []GetWalletResponse `json:"wallets"`
}

type UnlockWalletsRequest struct {
	Password string `json:"password"`
}
type UnlockWalletsResponse struct{}
