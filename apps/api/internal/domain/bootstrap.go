package domain

type GetBootstrapRequest struct{}
type GetBootstrapResponse struct {
	WalletCount   int
	WalletsLocked bool
	Settings      Settings
}
