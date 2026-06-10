package users_ws_transport

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:""`
}

type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}
