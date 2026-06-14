package friendships_service

import (
	"context"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type FriendshipsService struct {
	friendshipsRep    FriendshipsRepository
	friendRequestsRep FriendRequestsRepository
	broadcaster       FriendshipsWSTransport
}

type FriendshipsRepository interface {
	Save(ctx context.Context, requestID int) (domains.Friendship, error)
	FindByUserID(ctx context.Context, userID int, limit *int, offset *int) ([]domains.Friendship, error)
	Delete(ctx context.Context, userID int, friendshipID int) (domains.Friendship, error)
}

type FriendRequestsRepository interface {
	FindByIDAndToID(ctx context.Context, requestID int, toID int) (domains.FriendRequest, error)
}

type FriendshipsWSTransport interface {
	NotifyClientEvent(userID int, event string, friendship domains.Friendship)
	NotifyDeletedFriendship(userID int, friendshipID int)
}

func NewService(
	friendshipsRep FriendshipsRepository,
	friendRequestsRep FriendRequestsRepository,
	broadcaster FriendshipsWSTransport,
) *FriendshipsService {
	return &FriendshipsService{
		friendshipsRep:    friendshipsRep,
		friendRequestsRep: friendRequestsRep,
		broadcaster:       broadcaster,
	}
}

func (s *FriendshipsService) CreateFriendship(ctx context.Context, userID int, requestID int) (domains.Friendship, error) {
	if requestID <= 0 {
		return domains.Friendship{}, fmt.Errorf("request id must be positive: %w", core_errors.ErrInvalidArg)
	}

	friendRequest, err := s.friendRequestsRep.FindByIDAndToID(ctx, requestID, userID)
	if err != nil {
		return domains.Friendship{}, fmt.Errorf("failed to get friend request from rep: %w", err)
	}

	createdFriendship, err := s.friendshipsRep.Save(ctx, requestID)
	if err != nil {
		return domains.Friendship{}, fmt.Errorf("failed to create friendship: %w", err)
	}

	s.broadcaster.NotifyClientEvent(friendRequest.FromUser.ID, "accepted_friend_request", createdFriendship)
	s.broadcaster.NotifyClientEvent(friendRequest.ToUser.ID, "added_friend", createdFriendship)

	return createdFriendship, nil
}

func (s *FriendshipsService) GetFriendships(ctx context.Context, userID int, limit *int, offset *int) ([]domains.Friendship, error) {
	if limit != nil && *limit < 0 {
		return []domains.Friendship{}, fmt.Errorf("'limit' must be non negative: %w", core_errors.ErrInvalidArg)
	}

	if offset != nil && *offset < 0 {
		return []domains.Friendship{}, fmt.Errorf("'offset' must be non negative: %w", core_errors.ErrInvalidArg)
	}

	friendships, err := s.friendshipsRep.FindByUserID(ctx, userID, limit, offset)
	if err != nil {
		return []domains.Friendship{}, fmt.Errorf("failed to get friendships from rep: %w", err)
	}

	return friendships, nil
}

func (s *FriendshipsService) DeleteFriendship(ctx context.Context, userID int, friendshipID int) error {
	if friendshipID <= 0 {
		return fmt.Errorf("friendship id must be positive")
	}

	deletedFriendship, err := s.friendshipsRep.Delete(ctx, userID, friendshipID)

	if err != nil {
		return fmt.Errorf("failed to delete friendship: %w", err)
	}

	s.broadcaster.NotifyDeletedFriendship(userID, deletedFriendship.FirstUser.ID)
	s.broadcaster.NotifyDeletedFriendship(userID, deletedFriendship.SecondUser.ID)

	return nil
}
