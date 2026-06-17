package chats_ws_transport

import (
	"encoding/json"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
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

func (t *ChatsWSTransport) NotifyCreatedChat(userID int, chat domains.Chat) {
	content := ChatDTOResponse{
		ID: chat.ID,
		FirstUser: UserDTOResponse{
			ID:       chat.FirstUser.ID,
			Username: chat.FirstUser.Username,
			IsOnline: chat.FirstUser.IsOnline,
		},
		SecondUser: UserDTOResponse{
			ID:       chat.SecondUser.ID,
			Username: chat.SecondUser.Username,
			IsOnline: chat.SecondUser.IsOnline,
		},
		LastMessageContent: chat.LastMessageContent,
		LastMessageAt:      chat.LastMessageAt,
	}

	message := WebSocketMessage{
		Type:    "chat.created",
		Content: content,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create message to notify", zap.Error(err))

		return
	}

	fmt.Println(message)

	t.ws.NotifyClient(userID, msg)
}

func (t *ChatsWSTransport) NotifyDeletedChat(userID int, chatID int) {
	message := WebSocketMessage{
		Type:    "chat.deleted",
		Content: map[string]int{"chat_id": chatID},
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create message to notify", zap.Error(err))

		return
	}

	t.ws.NotifyClient(userID, msg)
}
