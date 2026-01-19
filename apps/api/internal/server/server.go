package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/codemaestro64/filament/apps/api/internal/service"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	httpServer *http.Server
	cfg        config.ServerConfig
	cancelApp  context.CancelFunc
}

func New(cfg config.ServerConfig, srvc *service.Service, cancelApp context.CancelFunc) (*Server, error) {
	mux := http.NewServeMux()

	registerHandlers(mux, srvc)

	//handler := recoveryMiddleware(mux)
	// H2C for HTTP/2 support without TLS
	h2s := &http2.Server{}
	h2cHandler := h2c.NewHandler(mux, h2s)

	return &Server{
		cfg:       cfg,
		cancelApp: cancelApp,
		httpServer: &http.Server{
			Handler:           h2cHandler,
			ReadHeaderTimeout: 2 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}, nil
}

func (s *Server) Name() string { return "api-server" }

func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("server failed to bind %s: %w", addr, err)
	}

	// Of configured port is 0, port was assigned dynamically
	if s.cfg.Port == 0 {
		s.cfg.Port = listener.Addr().(*net.TCPAddr).Port
	}

	go func() {
		log.Info().Int("port", s.cfg.Port).Msg("server listening")

		if err := s.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("critical server failure")

			// Invalidate the app context and triggers graceful shutdown of all components
			s.cancelApp()
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// --- Internal Helpers ---

func registerHandlers(mux *http.ServeMux, srvc *service.Service) {
	// Interceptors for ConnectRPC (Business Logic)
	/**opts := connect.WithInterceptors(
		interceptors.LoggingUnaryHandler(),
	)**/
}
