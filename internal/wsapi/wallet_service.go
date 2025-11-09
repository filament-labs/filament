package wsapi

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/filament-labs/filament/internal/pb"
	"github.com/filament-labs/filament/internal/pb/pbconnect"
	"github.com/filament-labs/filament/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletServer struct {
	walletService service.WalletService
}

func NewWalletServer(srvc *service.Service) (string, http.Handler) {
	walletServer := &WalletServer{
		walletService: srvc.Wallet,
	}

	return pbconnect.NewWalletServiceHandler(walletServer)
}

func (s *WalletServer) GetWallets(
	ctx context.Context,
	req *connect.Request[pb.WalletsRequest],
) (*connect.Response[pb.WalletsResponse], error) {

	wallets, err := s.walletService.GetWallets(ctx)
	if err != nil {
		return nil, err
	}

	resp := &pb.WalletsResponse{
		IsLocked: wallets.Locked,
	}
	for _, wal := range wallets.Wallets {
		resp.Wallets = append(resp.Wallets, &pb.WalletResponse{
			Id:         wal.ID,
			IsDefault:  wal.IsDefault,
			WalletName: wal.Name,
			Addresses:  wal.Addresses,
			Balance:    0,
			CreatedAt:  timestamppb.New(wal.CreatedAt),
		})
	}

	return connect.NewResponse(resp), nil
}
