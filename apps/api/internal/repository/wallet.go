package repository

import (
	"context"
	"fmt"

	"github.com/codemaestro64/filament/apps/api/internal/database/orm"
	"github.com/codemaestro64/filament/apps/api/internal/database/orm/wallet"
	"github.com/codemaestro64/filament/libs/filwallet"
	pbv1 "github.com/codemaestro64/filament/libs/proto/gen/go/v1"
)

type WalletRepo interface {
	CountWallets(ctx context.Context) (int, error)
	FindWallet(ctx context.Context, walletID int) (*filwallet.Wallet, error)
	GetWallets(ctx context.Context) ([]*filwallet.Wallet, error)
	DeleteWallet(ctx context.Context, walletID int) error
	SaveWallet(ctx context.Context, saveParams filwallet.SaveWalletParams) (*filwallet.Wallet, error)
}

type walletRepo struct {
	db *orm.Client
}

func newWalletRepo(db *orm.Client) WalletRepo {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) CountWallets(ctx context.Context) (int, error) {
	count, err := r.db.Wallet.Query().Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("db: count wallets: %w", err)
	}

	return count, nil
}

func (r *walletRepo) FindWallet(ctx context.Context, walletID int) (*filwallet.Wallet, error) {
	wallet, err := r.db.Wallet.Query().
		Where(wallet.ID(walletID)).
		WithAddresses().
		First(ctx)

	if err != nil && !orm.IsNotFound(err) {
		return nil, fmt.Errorf("db: find wallet by ID: %w", err)
	}

	wal := &filwallet.Wallet{
		IsDefault:        wallet.IsDefault,
		Name:             wallet.Name,
		EncryptedSeed:    wallet.EncryptedSeed,
		Salt:             wallet.Salt,
		EncryptedKeyJSON: wallet.EncryptedKeyJSON,
		CreatedAt:        wallet.CreatedAt,
	}

	for _, addr := range wallet.Edges.Addresses {
		wal.Addresses = append(wal.Addresses, &pbv1.Address{
			Type:  addr.Type,
			Value: addr.Address,
		})
	}

	return wal, nil
}

func (r *walletRepo) GetWallets(ctx context.Context) ([]*filwallet.Wallet, error) {
	return nil, nil
}

func (r *walletRepo) DeleteWallet(ctx context.Context, walletID int) error {
	err := r.db.Wallet.DeleteOneID(walletID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("db: delete wallet by ID: %w", err)
	}

	return nil
}

func (r *walletRepo) SaveWallet(ctx context.Context, saveParams filwallet.SaveWalletParams) (*filwallet.Wallet, error) {
	return nil, nil
}

func toAddressProto(addr *orm.Address) *pbv1.Address {
	if addr == nil {
		return nil
	}
	return &pbv1.Address{
		Type:  addr.Type,
		Value: addr.Address,
	}
}
