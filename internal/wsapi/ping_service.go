package wsapi

import (
	"context"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/filament-labs/filament/internal/pb"
	"github.com/filament-labs/filament/internal/pb/pbconnect"
)

type PingServer struct {
}

func NewPingServer() (string, http.Handler) {
	return pbconnect.NewPingServiceHandler(&PingServer{})
}

func (s *PingServer) Ping(
	ctx context.Context,
	req *connect.Request[pb.PingRequest],
) (*connect.Response[pb.PingResponse], error) {
	resp := &pb.PingResponse{
		Timestamp: time.Now().UnixMilli(),
	}

	return connect.NewResponse(resp), nil
}
