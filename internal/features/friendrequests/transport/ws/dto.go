package friendrequests_ws_transport

import (
	"time"
)

type FriendRequestDTOResponse struct {
	ID        int             `json:"request_id"`
	FromUser  UserDTOResponse `json:"from_user"`
	ToUser    UserDTOResponse `json:"to_user"`
	CreatedAt time.Time       `json:"created_at"`
}

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}
