package users_ws_transport

import (
	"encoding/json"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	"go.uber.org/zap"
)

type UsersWSTransport struct {
	ws *core_ws_server.Server
}

func NewWSTransport(ws *core_ws_server.Server) *UsersWSTransport {
	return &UsersWSTransport{
		ws: ws,
	}
}

func (t *UsersWSTransport) BroadcastEvent(event string, user domains.User) {
	content := UserDTO{
		ID:       user.ID,
		Username: user.Username,
		IsOnline: user.IsOnline,
	}

	message := Message{
		Type:    event,
		Content: content,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error(event, zap.Error(err))
		return
	}

	t.ws.Broadcast(msg)
}
