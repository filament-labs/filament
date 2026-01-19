package handler

import (
	"context"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/codemaestro64/filament/apps/api/internal/domain"
	"github.com/codemaestro64/filament/apps/api/internal/service"
	pbv1 "github.com/codemaestro64/filament/libs/proto/gen/go/v1"
	"github.com/codemaestro64/filament/libs/proto/gen/go/v1/pbv1connect"
)

type Request[T any] = connect.Request[T]
type Response[T any] = connect.Response[T]

type UserServer struct {
	userService service.UserService
}

func NewUserServer(srvc *service.Service, options connect.Option) (string, http.Handler) {
	userServer := &UserServer{
		userService: srvc.User,
	}

	return pbv1connect.NewUserServiceHandler(userServer, options)
}

func (s *UserServer) Bootstrap(
	ctx context.Context,
	_ *Request[pbv1.GetBootstrapRequest],
) (*Response[pbv1.GetBootstrapResponse], error) {

	result, err := s.userService.GetBootstrap(ctx, domain.GetBootstrapRequest{})
	if err != nil {
		return nil, nil
	}

	var network pbv1.NetworkType
	if result.Settings.Network.IsMainnet() {
		network = pbv1.NetworkType_NETWORK_MAINNET
	} else {
		network = pbv1.NetworkType_NETWORK_CALIBRATION_NET
	}

	resp := &pbv1.GetBootstrapResponse{
		WalletCount: int64(result.WalletCount),
		Settings: &pbv1.Settings{
			Network: network,
		},
	}

	return connect.NewResponse(resp), nil
}
