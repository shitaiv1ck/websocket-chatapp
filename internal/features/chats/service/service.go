package chats_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type ChatsService struct {
	chatsRep       ChatsRepository
	friendshipsRep FriendshipsRepository
	broadcaster    ChatsWSTransport
}

type ChatsRepository interface {
	SaveOrFind(ctx context.Context, firstUserID int, secondUserID int) (domains.Chat, error)
	FindByUserID(ctx context.Context, userID int, limit *int, offset *int) ([]domains.Chat, error)
	Delete(ctx context.Context, userID int, chatID int) (domains.Chat, error)
}

type FriendshipsRepository interface {
	FindByUsers(ctx context.Context, firstUserID int, secondUserID int) (domains.Friendship, error)
}

type ChatsWSTransport interface {
	NotifyDeletedChat(userID int, chatID int)
	NotifyCreatedChat(userID int, chat domains.Chat)
}

func NewService(
	chatsRep ChatsRepository,
	friendshipsRep FriendshipsRepository,
	broadcaster ChatsWSTransport,
) *ChatsService {
	return &ChatsService{
		chatsRep:       chatsRep,
		friendshipsRep: friendshipsRep,
		broadcaster:    broadcaster,
	}
}

func (s *ChatsService) CreateOrGetChat(ctx context.Context, userID int, friendID int) (domains.Chat, error) {
	if friendID <= 0 {
		return domains.Chat{}, fmt.Errorf("friend id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if userID == friendID {
		return domains.Chat{}, fmt.Errorf("can't create chat with yourself: %w", core_errors.ErrInvalidArg)
	}

	areFriends, err := s.areFriends(ctx, userID, friendID)
	if err != nil {
		return domains.Chat{}, fmt.Errorf("failed to get friendship: %w", err)
	}

	if !areFriends {
		return domains.Chat{}, fmt.Errorf("user with id=%v isn't your friend: %w", friendID, core_errors.ErrNotFound)
	}

	chat, err := s.chatsRep.SaveOrFind(ctx, userID, friendID)
	if err != nil {
		return domains.Chat{}, fmt.Errorf("failed to create or get chat: %w", err)
	}

	if chat.LastMessageContent == nil {
		s.broadcaster.NotifyCreatedChat(chat.FirstUser.ID, chat)
		s.broadcaster.NotifyCreatedChat(chat.SecondUser.ID, chat)
	}

	return chat, nil
}

func (s *ChatsService) areFriends(ctx context.Context, firstUserID int, secondUserID int) (bool, error) {
	_, err := s.friendshipsRep.FindByUsers(ctx, firstUserID, secondUserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *ChatsService) GetChats(ctx context.Context, userID int, limit *int, offset *int) ([]domains.Chat, error) {
	if limit != nil && *limit < 0 {
		return []domains.Chat{}, fmt.Errorf("'limit' must be non negative: %w", core_errors.ErrInvalidArg)
	}

	if offset != nil && *offset < 0 {
		return []domains.Chat{}, fmt.Errorf("'offset' must be non negative: %w", core_errors.ErrInvalidArg)
	}

	chats, err := s.chatsRep.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return []domains.Chat{}, fmt.Errorf("failed to get chats: %w", err)
	}

	return chats, nil
}

func (s *ChatsService) DeleteChat(ctx context.Context, userID int, chatID int) error {
	if chatID <= 0 {
		return fmt.Errorf("chat id must be positive: %w", core_errors.ErrInvalidArg)
	}

	deletedChat, err := s.chatsRep.Delete(ctx, userID, chatID)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	s.broadcaster.NotifyDeletedChat(deletedChat.FirstUser.ID, chatID)
	s.broadcaster.NotifyDeletedChat(deletedChat.SecondUser.ID, chatID)

	return nil
}
