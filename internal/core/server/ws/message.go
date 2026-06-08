package core_ws_server

import "encoding/json"

type Message struct {
	Type    string          `json:"type"`
	To      int             `json:"to"`
	Content json.RawMessage `json:"content"`
}

const (
	FriendRequestType = "friend_request"
	AcceptRequestType = "accept_request"
	SendMessageType   = "send_message"
)
