package friendships_http_transport

import "github.com/shitaiv1ck/realtime-chat/internal/core/domains"

type CreateFriendshipRequest struct {
	FriendRequestID int `json:"friend_request_id" validate:"required"`
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

func ToDTOResponse(friendships []domains.Friendship) []FriendshipDTOResponse {
	response := make([]FriendshipDTOResponse, len(friendships))

	for i, friendship := range friendships {
		response[i] = FriendshipDTOResponse{
			ID: friendship.ID,
			FirstUser: UserDTOResponse{
				ID:       friendship.FirstUser.ID,
				Username: friendship.FirstUser.Username,
				IsOnline: friendship.FirstUser.IsOnline,
			},
			SecondUser: UserDTOResponse{
				ID:       friendship.SecondUser.ID,
				Username: friendship.SecondUser.Username,
				IsOnline: friendship.SecondUser.IsOnline,
			},
		}
	}

	return response
}
