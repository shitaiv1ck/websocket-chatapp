package friendrequests_ws_transport

import "time"

type FriendRequestDTO struct {
	ID        int       `json:"request_id"`
	FromUser  UserDTO   `json:"from_user"`
	ToUser    UserDTO   `json:"to_user"`
	CreatedAt time.Time `json:"created_at"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Message struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}
