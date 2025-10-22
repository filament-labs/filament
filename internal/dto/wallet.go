package dto

type GetWalletRequest struct {
	WalletID   string `json:"wallet_id"`
	WalletName string `json:"wallet_name"`
}

type GetWalletResponse struct {
	ID        string             `json:"id"`         // Unique wallet identifier
	IsDefault bool               `json:"is_default"` // Is user's default wallet
	Name      string             `json:"name"`       // User-provided wallet name
	Address   string             `json:"address"`    // Filecoin f1 address
	Balance   GetBalanceResponse `json:"balance"`    // Wallet balance
	CreatedAt int64              `json:"created_at"` // Timestamp
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
