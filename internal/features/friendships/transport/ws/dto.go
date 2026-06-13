package friendships_ws_transport

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type FriendshipDTOResponse struct {
	ID         int             `json:"id"`
	FirstUser  UserDTOResponse `json:"first_user"`
	SecondUser UserDTOResponse `json:"second_user"`
}

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}
