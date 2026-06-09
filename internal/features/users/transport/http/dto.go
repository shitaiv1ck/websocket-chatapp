package users_http_transport

import "github.com/shitaiv1ck/realtime-chat/internal/core/domains"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type CreateUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type GetUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type PatchUserRequest struct {
	Username domains.Nullable[string] `json:"username"`
	IsOnline domains.Nullable[bool]   `json:"is_online"`
}

type PatchUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

func usersToResponse(users []domains.User) []GetUserResponse {
	response := make([]GetUserResponse, 0)

	for _, user := range users {
		userToResponse := GetUserResponse{
			ID:       user.ID,
			Username: user.Username,
			IsOnline: user.IsOnline,
		}

		response = append(response, userToResponse)
	}

	return response
}
