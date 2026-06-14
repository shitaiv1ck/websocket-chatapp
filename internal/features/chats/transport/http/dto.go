package chats_http_transport

import (
	"time"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
)

type CreateOrGetChatRequest struct {
	FriendID int `json:"friend_id" validate:"required"`
}

type UserDTOResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type ChatDTOResponse struct {
	ID                 int             `json:"id"`
	FirstUser          UserDTOResponse `json:"first_user"`
	SecondUser         UserDTOResponse `json:"second_user"`
	LastMessageContent *string         `json:"last_message_content"`
	LastMessageAt      time.Time       `json:"last_message_at"`
}

func ToDTOResponse(chats []domains.Chat) []ChatDTOResponse {
	response := make([]ChatDTOResponse, len(chats))

	for i, chat := range chats {
		response[i] = ChatDTOResponse{
			ID: chat.ID,
			FirstUser: UserDTOResponse{
				ID:       chat.FirstUser.ID,
				Username: chat.FirstUser.Username,
				IsOnline: chat.FirstUser.IsOnline,
			},
			SecondUser: UserDTOResponse{
				ID:       chat.SecondUser.ID,
				Username: chat.SecondUser.Username,
				IsOnline: chat.SecondUser.IsOnline,
			},
			LastMessageContent: chat.LastMessageContent,
			LastMessageAt:      chat.LastMessageAt,
		}
	}

	return response
}
