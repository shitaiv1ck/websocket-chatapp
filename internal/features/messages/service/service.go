package messages_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type MessagesService struct {
	messageRep     MessagesRepository
	friendshipsRep FriendshipsRepository
	broadcaster    MessagesWSTransport
}

type MessagesRepository interface {
	Save(ctx context.Context, message domains.Message) (domains.Message, error)
	FindByChatIDAndUserID(ctx context.Context, userID, chatID int) ([]domains.Message, error)
}

type FriendshipsRepository interface {
	FindByUsers(ctx context.Context, firstUserID int, secondUserID int) (domains.Friendship, error)
}

type MessagesWSTransport interface {
	NotifyClientEvent(userID int, event string, message domains.Message)
}

func NewService(
	messageRep MessagesRepository,
	friendshipsRep FriendshipsRepository,
	broadcaster MessagesWSTransport,
) *MessagesService {
	return &MessagesService{
		messageRep:     messageRep,
		friendshipsRep: friendshipsRep,
		broadcaster:    broadcaster,
	}
}

func (s *MessagesService) CreateMessage(ctx context.Context, message domains.Message) (domains.Message, error) {
	if err := message.Validate(); err != nil {
		return domains.Message{}, fmt.Errorf("failed to validate message: %w", err)
	}

	areFriends, err := s.areFriends(ctx, message.SenderID, message.ReceiverID)
	if err != nil {
		return domains.Message{}, fmt.Errorf("failed to get friendships: %w", err)
	}

	if !areFriends {
		return domains.Message{}, fmt.Errorf("user with id=%v isn't your friend: %w", message.ReceiverID, core_errors.ErrInvalidArg)
	}

	createdMessage, err := s.messageRep.Save(ctx, message)
	if err != nil {
		return domains.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	s.broadcaster.NotifyClientEvent(createdMessage.SenderID, "message.sent", createdMessage)
	s.broadcaster.NotifyClientEvent(createdMessage.ReceiverID, "message.received", createdMessage)

	return createdMessage, nil
}

func (s *MessagesService) areFriends(ctx context.Context, firstUserID int, secondUserID int) (bool, error) {
	_, err := s.friendshipsRep.FindByUsers(ctx, firstUserID, secondUserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *MessagesService) GetMessages(ctx context.Context, userID int, chatID int) ([]domains.Message, error) {
	if chatID <= 0 {
		return []domains.Message{}, fmt.Errorf("chat id must be positive: %w", core_errors.ErrInvalidArg)
	}

	messages, err := s.messageRep.FindByChatIDAndUserID(ctx, userID, chatID)
	if err != nil {
		return []domains.Message{}, fmt.Errorf("failed to get messages: %w", err)
	}

	return messages, nil
}
