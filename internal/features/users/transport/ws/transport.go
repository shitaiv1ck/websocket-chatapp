package users_ws_transport

import (
	"encoding/json"

	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
)

type UsersWSTransport struct {
	ws *core_ws_server.Server
}

func NewWSTransport(ws *core_ws_server.Server) *UsersWSTransport {
	return &UsersWSTransport{
		ws: ws,
	}
}

func (t *UsersWSTransport) BroadcastEvent(event string, content any) error {
	message := Message{
		Type:    event,
		Content: content,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	t.ws.Broadcast(msg)

	return nil
}
