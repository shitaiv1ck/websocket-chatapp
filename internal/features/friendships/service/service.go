package friendships_service

import (
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
	Save(fromUserID int, toUserID int) (domains.Friendship, error)
	FindByUserID(userID int, limit *int, offset *int) ([]domains.Friendship, error)
	Delete(userID int, friendshipID int) error
}

type FriendRequestsRepository interface {
	FindByIDAndToID(requestID int, toID int) (domains.FriendRequest, error)
	Delete(userID int, requestID int) error
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

func (s *FriendshipsService) CreateFriendship(userID int, requestID int) (domains.Friendship, error) {
	if requestID <= 0 {
		return domains.Friendship{}, fmt.Errorf("request id must be positive: %w", core_errors.ErrInvalidArg)
	}

	friendRequest, err := s.friendRequestsRep.FindByIDAndToID(requestID, userID)
	if err != nil {
		return domains.Friendship{}, fmt.Errorf("failed to get friend request from rep: %w", err)
	}

	if err := s.friendRequestsRep.Delete(userID, requestID); err != nil {
		return domains.Friendship{}, fmt.Errorf("failed to delete friend request: %w", err)
	}

	createdFriendship, err := s.friendshipsRep.Save(friendRequest.FromUser.ID, friendRequest.ToUser.ID)
	if err != nil {
		return domains.Friendship{}, fmt.Errorf("failed to create friendship: %w", err)
	}

	s.broadcaster.NotifyClientEvent(friendRequest.FromUser.ID, "accepted_friend_request", createdFriendship)
	s.broadcaster.NotifyClientEvent(friendRequest.ToUser.ID, "friend_added", createdFriendship)

	return createdFriendship, nil
}

func (s *FriendshipsService) GetFriendships(userID int, limit *int, offset *int) ([]domains.Friendship, error) {
	if limit != nil && *limit < 0 {
		return []domains.Friendship{}, fmt.Errorf("'limit' must be non negative: %w", core_errors.ErrInvalidArg)
	}

	if offset != nil && *offset < 0 {
		return []domains.Friendship{}, fmt.Errorf("'offset' must be non negative: %w", core_errors.ErrInvalidArg)
	}

	friendships, err := s.friendshipsRep.FindByUserID(userID, limit, offset)
	if err != nil {
		return []domains.Friendship{}, fmt.Errorf("failed to get friendships from rep: %w", err)
	}

	return friendships, nil
}

func (s *FriendshipsService) DeleteFriendship(userID int, friendshipID int) error {
	if friendshipID <= 0 {
		return fmt.Errorf("friendship id must be positive")
	}

	if err := s.friendshipsRep.Delete(userID, friendshipID); err != nil {
		return fmt.Errorf("failed to delete friendship: %w", err)
	}

	s.broadcaster.NotifyDeletedFriendship(userID, friendshipID)

	return nil
}
