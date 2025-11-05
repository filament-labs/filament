package wsapi

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"

	"github.com/filament-labs/filament/internal/app"
	"github.com/filament-labs/filament/internal/service"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type HandlerFunc func(*centrifuge.Client, json.RawMessage) error

type WsAPI struct {
	node               *centrifuge.Node
	handlers           map[string]HandlerFunc
	transactionService service.TransactionService
}

func New() app.App {
	return &WsAPI{}
}

func (ws *WsAPI) Bootstrap(srvc *service.Service) error {
	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		return fmt.Errorf("error initializing websocket server: %w", err)
	}

	ws.node = node
	ws.handlers = make(map[string]HandlerFunc)
	ws.registerHandlers()

	return nil
}

func (ws *WsAPI) registerHandlers() {

}

func (ws *WsAPI) Run() error {
	ws.node.OnConnect(func(client *centrifuge.Client) {
		// start transaction listener

		client.OnMessage(func(e centrifuge.MessageEvent) {
			var msg Message
			if err := json.Unmarshal(e.Data, &msg); err != nil {
				//return centrifuge.MessageReply{Error: centrifuge.ErrorBadRequest}
			}

			handler, exists := ws.handlers[msg.Type]
			if !exists {
				/**return centrifuge.MessageReply{
					Error: centrifuge.Error{Code: 400, Message: "unknown message type"},
				}**/
			}

			if err := handler(client, msg.Payload); err != nil {
				/**return centrifuge.MessageReply{
					Error: centrifuge.Error{Code: 500, Message: err.Error()},
				}**/
			}

		})
	})

	err := ws.node.Run()
	if err != nil {
		return fmt.Errorf("error running websocker server node: %w", err)
	}

	return nil
}
