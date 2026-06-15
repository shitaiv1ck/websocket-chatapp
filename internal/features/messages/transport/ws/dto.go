package messages_ws_transport

import "time"

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type MessageDTOResponse struct {
	ID         int       `json:"id"`
	ChatID     int       `json:"chat_id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
