package chats_ws_transport

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}
