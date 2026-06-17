package chats_ws_transport

import "time"

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type ChatDTOResponse struct {
	ID                 int             `json:"id"`
	FirstUser          UserDTOResponse `json:"first_user"`
	SecondUser         UserDTOResponse `json:"second_user"`
	LastMessageContent *string         `json:"last_message_content"`
	LastMessageAt      time.Time       `json:"last_message_at"`
}
