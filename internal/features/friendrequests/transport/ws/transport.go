package friendrequests_ws_transport

import (
	"encoding/json"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_ws_server "github.com/shitaiv1ck/realtime-chat/internal/core/server/ws"
	"go.uber.org/zap"
)

type FriendRequestsWSTransport struct {
	ws *core_ws_server.Server
}

func NewWSTransport(ws *core_ws_server.Server) *FriendRequestsWSTransport {
	return &FriendRequestsWSTransport{
		ws: ws,
	}
}

func (t *FriendRequestsWSTransport) NotifyClientEvent(userID int, event string, request domains.FriendRequest) {
	content := FriendRequestDTO{
		ID:        request.ID,
		FromUser:  UserDTO{ID: request.FromUser.ID, Username: request.FromUser.Username},
		ToUser:    UserDTO{ID: request.ToUser.ID, Username: request.ToUser.Username},
		CreatedAt: request.CreatedAt,
	}

	message := Message{
		Type:    event,
		Content: content,
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create message to notify", zap.Error(err))

		return
	}

	t.ws.NotifyClient(userID, msg)
}

func (t *FriendRequestsWSTransport) NotifyDeclinedRequest(userID int, requestID int) {
	message := Message{
		Type:    "declined_friend_request",
		Content: map[string]int{"request_id": requestID},
	}

	msg, err := json.Marshal(message)
	if err != nil {
		t.ws.GetLogger().Error("failed to create message to notify", zap.Error(err))

		return
	}

	t.ws.NotifyClient(userID, msg)
}
