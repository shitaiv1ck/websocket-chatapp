package users_http_transport

import (
	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
)

type UserDTORequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type PatchUserRequest struct {
	Username    domains.Nullable[string] `json:"username"`
	OldPassword domains.Nullable[string] `json:"old_password"`
	NewPassword domains.Nullable[string] `json:"new_password"`
}

func ToDTOResponse(users []domains.User) []UserDTOResponse {
	response := make([]UserDTOResponse, len(users))

	for i, user := range users {
		response[i] = UserDTOResponse{
			ID:       user.ID,
			Username: user.Username,
			IsOnline: user.IsOnline,
		}
	}

	return response
}
