package friendrequests_http_transport

import (
	"time"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
)

type CreateFriendRequestRequest struct {
	ToUserID int `json:"to_user_id" validate:"required"`
}

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type FriendRequestDTOResponse struct {
	ID        int             `json:"id"`
	FromUser  UserDTOResponse `json:"from_user"`
	ToUser    UserDTOResponse `json:"to_user"`
	CreatedAt time.Time       `json:"created_at"`
}

func ToDTOResponse(requests []domains.FriendRequest) []FriendRequestDTOResponse {
	response := make([]FriendRequestDTOResponse, len(requests))

	for i, request := range requests {
		response[i] = FriendRequestDTOResponse{
			ID: request.ID,
			FromUser: UserDTOResponse{
				ID:       request.FromUser.ID,
				Username: request.FromUser.Username,
				IsOnline: request.FromUser.IsOnline,
			},
			ToUser: UserDTOResponse{
				ID:       request.ToUser.ID,
				Username: request.ToUser.Username,
				IsOnline: request.ToUser.IsOnline,
			},
			CreatedAt: request.CreatedAt,
		}
	}

	return response
}
