package friendrequests_service

import (
	"errors"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type FriendRequestsService struct {
	friendRequestsRep FriendRequestsRepository
	friendshipRep     FriendshipRepository
	broadcaster       FriendRequestsWSTransport
}

type FriendRequestsRepository interface {
	Save(request domains.FriendRequest) (domains.FriendRequest, error)
	FindByFromIDAndToID(fromID int, toID int) (domains.FriendRequest, error)
	FindByToID(toID int) ([]domains.FriendRequest, error)
	FindByFromID(fromID int) ([]domains.FriendRequest, error)
	Delete(userID int, requestID int) error
}

type FriendshipRepository interface {
	FindByUsers(firstUserID int, secondUserID int) (domains.Friendship, error)
}

type FriendRequestsWSTransport interface {
	NotifyClientEvent(userID int, event string, request domains.FriendRequest)
	NotifyDeclinedRequest(userID int, requestID int)
}

func NewService(
	friendRequestsRep FriendRequestsRepository,
	friendshipRep FriendshipRepository,
	broadcaster FriendRequestsWSTransport,
) *FriendRequestsService {
	return &FriendRequestsService{
		friendRequestsRep: friendRequestsRep,
		friendshipRep:     friendshipRep,
		broadcaster:       broadcaster,
	}
}

func (s *FriendRequestsService) CreateFriendRequest(request domains.FriendRequest) (domains.FriendRequest, error) {
	if err := request.Validate(); err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to validate friend request: %w", err)
	}

	hasIncoming, err := s.hasIncomingFromTargetUser(request.ToUser.ID, request.FromUser.ID)
	if err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to get friend request: %w", err)
	}

	if hasIncoming {
		return domains.FriendRequest{}, fmt.Errorf("user with id=%v already sent friend request: %w", request.ToUser.ID, core_errors.ErrConflict)
	}

	areFriends, err := s.areFriends(request.FromUser.ID, request.ToUser.ID)
	if err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to get friendship: %w", err)
	}

	if areFriends {
		return domains.FriendRequest{}, fmt.Errorf("user with id=%v already your friend: %w", request.ToUser.ID, core_errors.ErrConflict)
	}

	createdFriendRequest, err := s.friendRequestsRep.Save(request)
	if err != nil {
		return domains.FriendRequest{}, err
	}

	s.broadcaster.NotifyClientEvent(request.ToUser.ID, "received_friend_request", createdFriendRequest)
	s.broadcaster.NotifyClientEvent(request.FromUser.ID, "sent_friend_request", createdFriendRequest)

	return createdFriendRequest, nil
}

func (s *FriendRequestsService) hasIncomingFromTargetUser(toUserID int, fromUserID int) (bool, error) {
	_, err := s.friendRequestsRep.FindByFromIDAndToID(toUserID, fromUserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *FriendRequestsService) areFriends(firstUserID int, secondUserID int) (bool, error) {
	_, err := s.friendshipRep.FindByUsers(firstUserID, secondUserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *FriendRequestsService) GetFriendRequests(userID int, direction *string) ([]domains.FriendRequest, error) {
	if direction == nil {
		in := "incoming"
		direction = &in
	}

	requests := make([]domains.FriendRequest, 0)
	var err error

	switch *direction {
	case "incoming":
		requests, err = s.friendRequestsRep.FindByToID(userID)
	case "outgoing":
		requests, err = s.friendRequestsRep.FindByFromID(userID)
	default:
		return []domains.FriendRequest{}, fmt.Errorf("invalid direction '%s': %w", *direction, core_errors.ErrInvalidArg)
	}

	if err != nil {
		return []domains.FriendRequest{}, fmt.Errorf("failed to get friend requests: %w", err)
	}

	return requests, nil
}

func (s *FriendRequestsService) DeleteFriendRequest(userID int, requestID int) error {
	if err := s.friendRequestsRep.Delete(userID, requestID); err != nil {
		return err
	}

	s.broadcaster.NotifyDeclinedRequest(userID, requestID)

	return nil
}
