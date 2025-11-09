package wsapi

import (
	"net/http"

	"github.com/filament-labs/filament/internal/app"
	"github.com/filament-labs/filament/internal/service"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"

	connectcors "connectrpc.com/cors"
)

type WsAPI struct {
	srvc *service.Service
	mux  *http.ServeMux
}

func New() app.App {
	return &WsAPI{}
}

func (ws *WsAPI) Bootstrap(srvc *service.Service) error {
	ws.srvc = srvc
	ws.mux = http.NewServeMux()
	ws.registerServices()
	return nil
}

func (ws *WsAPI) registerServices() {
	log.Info().Msg("registering service handlers")
	pingPath, pingHandler := NewPingServer()
	walletPath, walletHandler := NewWalletServer(ws.srvc)

	ws.mux.Handle(pingPath, pingHandler)
	ws.mux.Handle(walletPath, walletHandler)
}

func (ws *WsAPI) Run() error {
	log.Info().Str("port", "8080").Msg("Starting server")
	handler := withCORS(ws.mux)

	return http.ListenAndServe(":8080", handler)
}

func withCORS(connectHandler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
		MaxAge:         7200,
	})
	return c.Handler(connectHandler)
}
