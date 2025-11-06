package wsapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/centrifugal/centrifuge"
	"github.com/filament-labs/filament/internal/app"
	"github.com/filament-labs/filament/internal/service"
	"github.com/rs/zerolog/log"
)

// Message represents the wire format for client messages.
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// HandlerFunc processes a message payload for a connected client.
type HandlerFunc func(*centrifuge.Client, json.RawMessage) ([]byte, error)

type WsAPI struct {
	node     *centrifuge.Node
	handlers map[string]HandlerFunc
}

// New creates a new WsAPI instance.
func New() app.App {
	return &WsAPI{
		handlers: make(map[string]HandlerFunc),
	}
}

// Bootstrap initializes the Centrifuge node and registers message handlers.
func (ws *WsAPI) Bootstrap(srvc *service.Service) error {
	log.Debug().Msg("Initializing Centrifuge WebSocket server")

	node, err := centrifuge.New(centrifuge.Config{
		LogLevel: centrifuge.LogLevelError,
	})
	if err != nil {
		return fmt.Errorf("failed to create Centrifuge node: %w", err)
	}

	ws.node = node
	ws.registerHandlers()

	log.Info().Msg("WebSocket API handlers registered")
	return nil
}

// registerHandlers sets up all supported message types.
func (ws *WsAPI) registerHandlers() {
	ws.handlers["ping"] = ws.handlePing
	ws.handlers["subscribe"] = ws.handleSubscribe
}

// Run starts the Centrifuge server and begins accepting connections.
func (ws *WsAPI) Run() error {
	ws.node.OnConnect(func(client *centrifuge.Client) {
		log.Info().Str("client", client.ID()).Msg("Client connected")

		client.OnMessage(func(e centrifuge.MessageEvent) {
			reply := ws.handleClientMessage(client, e.Data)
			if reply != nil {
				_ = client.Send(reply)
			}
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			log.Info().Str("client", client.ID()).Msg("Client disconnected")
		})
	})

	log.Info().Msg("Starting WebSocket server")
	if err := ws.node.Run(); err != nil {
		return fmt.Errorf("failed to run Centrifuge server: %w", err)
	}

	wsHandler := centrifuge.NewWebsocketHandler(ws.node, centrifuge.WebsocketConfig{})
	http.Handle("/connection/websocket", wsHandler)
	log.Info().Msg("Starting server, visit http://localhost:8000")

	return http.ListenAndServe(":8000", nil)
}

func (ws *WsAPI) handleClientMessage(client *centrifuge.Client, data []byte) []byte {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Warn().Str("client", client.ID()).Err(err).Msg("Invalid message format")
		return ws.errorReply(400, "invalid JSON")
	}

	handler, exists := ws.handlers[msg.Type]
	if !exists {
		log.Warn().Str("client", client.ID()).Str("type", msg.Type).Msg("Unknown message type")
		return ws.errorReply(400, "unknown message type")
	}

	resp, err := handler(client, msg.Payload)
	if err != nil {
		log.Error().Str("client", client.ID()).Str("type", msg.Type).Err(err).Msg("Handler error")
		return ws.errorReply(500, err.Error())
	}

	return resp
}

func (ws *WsAPI) errorReply(code int, message string) []byte {
	resp := map[string]any{
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	}
	data, _ := json.Marshal(resp)
	return data
}

func (ws *WsAPI) handlePing(_ *centrifuge.Client, _ json.RawMessage) ([]byte, error) {
	resp := map[string]string{"type": "pong"}
	return json.Marshal(resp)
}

func (ws *WsAPI) handleSubscribe(client *centrifuge.Client, payload json.RawMessage) ([]byte, error) {
	var req struct {
		Channel string `json:"channel"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, fmt.Errorf("invalid subscribe payload: %w", err)
	}

	if req.Channel == "" {
		return nil, fmt.Errorf("channel name is required")
	}

	if err := client.Subscribe(req.Channel); err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	log.Info().Str("client", client.ID()).Str("channel", req.Channel).Msg("Subscribed to channel")

	resp := map[string]any{
		"type":    "subscribed",
		"channel": req.Channel,
	}
	return json.Marshal(resp)
}
