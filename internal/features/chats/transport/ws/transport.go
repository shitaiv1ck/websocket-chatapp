package chats_ws_transport

import (
	"encoding/json"

	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	"go.uber.org/zap"
)

type ChatsWSTransport struct {
	ws core_ws_server.Broadcaster
}

func NewWSTransport(ws core_ws_server.Broadcaster) *ChatsWSTransport {
	return &ChatsWSTransport{
		ws: ws,
	}
}

func (t *ChatsWSTransport) NotifyDeleteChat(userID int, chatID int) {
	message := WebSocketMessage{
		Type:    "deleted_chat",
		Content: map[string]int{"chat_id": chatID},
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create message to notify", zap.Error(err))

		return
	}

	t.ws.NotifyClient(userID, msg)
}
