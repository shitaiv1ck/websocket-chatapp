package users_ws_transport

import (
	"encoding/json"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	"go.uber.org/zap"
)

type UsersWSTransport struct {
	ws core_ws_server.Broadcaster
}

func NewWSTransport(ws core_ws_server.Broadcaster) *UsersWSTransport {
	return &UsersWSTransport{
		ws: ws,
	}
}

func (t *UsersWSTransport) BroadcastEvent(event string, user domains.User) {
	content := UserDTOResponse{
		ID:       user.ID,
		Username: user.Username,
		IsOnline: user.IsOnline,
	}

	message := WebSocketMessage{
		Type:    event,
		Content: content,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create msg to broadcast", zap.Error(err))
		return
	}

	t.ws.Broadcast(msg)
}
