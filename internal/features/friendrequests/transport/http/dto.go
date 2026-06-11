package friendrequests_http_transport

import (
	"time"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
)

type CreateFriendRequestRequest struct {
	ToUserID int `json:"to_user_id" validate:"required"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type CreateFriendRequestResponse struct {
	ID        int       `json:"id"`
	FromUser  UserDTO   `json:"from_user"`
	ToUser    UserDTO   `json:"to_user"`
	CreatedAt time.Time `json:"created_at"`
}

type GetFriendRequestResponse struct {
	ID        int       `json:"id"`
	FromUser  UserDTO   `json:"from_user"`
	ToUser    UserDTO   `json:"to_user"`
	CreatedAt time.Time `json:"created_at"`
}

func friendRequestsToResponse(requests []domains.FriendRequest) []GetFriendRequestResponse {
	response := make([]GetFriendRequestResponse, 0)

	for _, request := range requests {
		requestToResponse := GetFriendRequestResponse{
			ID:        request.ID,
			FromUser:  UserDTO{ID: request.FromUser.ID, Username: request.FromUser.Username},
			ToUser:    UserDTO{ID: request.ToUser.ID, Username: request.ToUser.Username},
			CreatedAt: request.CreatedAt,
		}

		response = append(response, requestToResponse)
	}

	return response
}
