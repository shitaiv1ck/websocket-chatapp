package messages_ws_transport

import (
	"encoding/json"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	"go.uber.org/zap"
)

type MessagesWSTransport struct {
	ws core_ws_server.Broadcaster
}

func NewWSTransport(ws core_ws_server.Broadcaster) *MessagesWSTransport {
	return &MessagesWSTransport{
		ws: ws,
	}
}

func (t *MessagesWSTransport) NotifyClientEvent(userID int, event string, message domains.Message) {
	content := MessageDTOResponse{
		ID:         message.ID,
		ChatID:     message.ChatID,
		SenderID:   message.SenderID,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt,
	}

	messageWS := WebSocketMessage{
		Type:    event,
		Content: content,
	}

	msg, err := json.Marshal(messageWS)
	if err != nil {
		t.ws.GetLogger().Error("failed to create msg to broadcast", zap.Error(err))
		return
	}

	t.ws.NotifyClient(userID, msg)
}
