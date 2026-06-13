package friendships_ws_transport

import (
	"encoding/json"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	"go.uber.org/zap"
)

type FriendshipsWSTransport struct {
	ws *core_ws_server.Server
}

func NewWSTransport(ws *core_ws_server.Server) *FriendshipsWSTransport {
	return &FriendshipsWSTransport{
		ws: ws,
	}
}

func (t *FriendshipsWSTransport) NotifyClientEvent(userID int, event string, friendship domains.Friendship) {
	content := FriendshipDTOResponse{
		ID: friendship.ID,
		FirstUser: UserDTOResponse{
			ID:       friendship.FirstUser.ID,
			Username: friendship.FirstUser.Username,
			IsOnline: friendship.FirstUser.IsOnline,
		},
		SecondUser: UserDTOResponse{
			ID:       friendship.SecondUser.ID,
			Username: friendship.SecondUser.Username,
			IsOnline: friendship.SecondUser.IsOnline,
		},
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

	t.ws.NotifyClient(userID, msg)
}

func (t *FriendshipsWSTransport) NotifyDeletedFriendship(userID int, friendshipID int) {
	message := WebSocketMessage{
		Type:    "deleted_friendship",
		Content: map[string]int{"friendship_id": friendshipID},
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create msg to broadcast", zap.Error(err))
		return
	}

	t.ws.NotifyClient(userID, msg)
}
