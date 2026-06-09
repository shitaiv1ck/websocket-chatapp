package users_ws_transport

type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}
