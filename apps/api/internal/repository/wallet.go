package repository

import (
	"context"
	"fmt"

	"github.com/codemaestro64/filament/apps/api/internal/infra/database/orm"
	dbwallet "github.com/codemaestro64/filament/apps/api/internal/infra/database/orm/wallet"
	"github.com/codemaestro64/filament/libs/filwallet"
	"github.com/codemaestro64/filament/libs/filwallet/address"
	"github.com/codemaestro64/filament/libs/filwallet/wallet"
)

type WalletRepo interface {
	CountWallets(ctx context.Context) (int, error)
	FindWallet(ctx context.Context, walletID int) (*wallet.Wallet, error)
	GetWallets(ctx context.Context) ([]*wallet.Wallet, error)
	DeleteWallet(ctx context.Context, walletID int) error
	SaveWallet(ctx context.Context, saveParams filwallet.SaveWalletParams) (*wallet.Wallet, error)
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

func (r *walletRepo) FindWallet(ctx context.Context, walletID int) (*wallet.Wallet, error) {
	dbWallet, err := r.db.Wallet.Query().
		Where(dbwallet.IDEQ(walletID)).
		WithAddresses().
		First(ctx)

	if err != nil && !orm.IsNotFound(err) {
		return nil, fmt.Errorf("db: find wallet by ID: %w", err)
	}

	wal := &wallet.Wallet{
		IsDefault:         dbWallet.IsDefault,
		Name:              dbWallet.Name,
		EncryptedMnemonic: dbWallet.EncryptedSeed,
		Salt:              dbWallet.Salt,
		EncryptedKeyJSON:  dbWallet.EncryptedKeyJSON,
		CreatedAt:         dbWallet.CreatedAt,
	}

	for _, addr := range dbWallet.Edges.Addresses {
		wal.Addresses = append(wal.Addresses, address.Address{
			Type:  addr.Type,
			Value: addr.Address,
		})
	}

	return wal, nil
}

func (r *walletRepo) GetWallets(ctx context.Context) ([]*wallet.Wallet, error) {
	return nil, nil
}

func (r *walletRepo) DeleteWallet(ctx context.Context, walletID int) error {
	err := r.db.Wallet.DeleteOneID(walletID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("db: delete wallet by ID: %w", err)
	}

	return nil
}

func (r *walletRepo) SaveWallet(ctx context.Context, saveParams filwallet.SaveWalletParams) (*wallet.Wallet, error) {
	return nil, nil
}
