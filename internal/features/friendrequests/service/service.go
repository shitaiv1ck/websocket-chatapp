package friendrequests_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type FriendRequestsService struct {
	friendRequestsRep FriendRequestsRepository
	friendshipsRep    FriendshipRepository
	broadcaster       FriendRequestsWSTransport
}

type FriendRequestsRepository interface {
	Save(ctx context.Context, request domains.FriendRequest) (domains.FriendRequest, error)
	FindByFromIDAndToID(ctx context.Context, fromID int, toID int) (domains.FriendRequest, error)
	FindByToID(ctx context.Context, toID int) ([]domains.FriendRequest, error)
	FindByFromID(ctx context.Context, fromID int) ([]domains.FriendRequest, error)
	Delete(ctx context.Context, userID int, requestID int) (domains.FriendRequest, error)
}

type FriendshipRepository interface {
	FindByUsers(ctx context.Context, firstUserID int, secondUserID int) (domains.Friendship, error)
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
		friendshipsRep:    friendshipRep,
		broadcaster:       broadcaster,
	}
}

func (s *FriendRequestsService) CreateFriendRequest(ctx context.Context, request domains.FriendRequest) (domains.FriendRequest, error) {
	if err := request.Validate(); err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to validate friend request: %w", err)
	}

	hasIncoming, err := s.hasIncomingFromTargetUser(ctx, request.ToUser.ID, request.FromUser.ID)
	if err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to get friend request: %w", err)
	}

	if hasIncoming {
		return domains.FriendRequest{}, fmt.Errorf("user with id=%v already sent friend request: %w", request.ToUser.ID, core_errors.ErrConflict)
	}

	areFriends, err := s.areFriends(ctx, request.FromUser.ID, request.ToUser.ID)
	if err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to get friendship: %w", err)
	}

	if areFriends {
		return domains.FriendRequest{}, fmt.Errorf("user with id=%v already your friend: %w", request.ToUser.ID, core_errors.ErrConflict)
	}

	createdFriendRequest, err := s.friendRequestsRep.Save(ctx, request)
	if err != nil {
		return domains.FriendRequest{}, err
	}

	s.broadcaster.NotifyClientEvent(request.ToUser.ID, "received_friend_request", createdFriendRequest)
	s.broadcaster.NotifyClientEvent(request.FromUser.ID, "sent_friend_request", createdFriendRequest)

	return createdFriendRequest, nil
}

func (s *FriendRequestsService) hasIncomingFromTargetUser(ctx context.Context, toUserID int, fromUserID int) (bool, error) {
	_, err := s.friendRequestsRep.FindByFromIDAndToID(ctx, toUserID, fromUserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *FriendRequestsService) areFriends(ctx context.Context, firstUserID int, secondUserID int) (bool, error) {
	_, err := s.friendshipsRep.FindByUsers(ctx, firstUserID, secondUserID)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *FriendRequestsService) GetFriendRequests(ctx context.Context, userID int, direction *string) ([]domains.FriendRequest, error) {
	if direction == nil {
		in := "incoming"
		direction = &in
	}

	requests := make([]domains.FriendRequest, 0)
	var err error

	switch *direction {
	case "incoming":
		requests, err = s.friendRequestsRep.FindByToID(ctx, userID)
	case "outgoing":
		requests, err = s.friendRequestsRep.FindByFromID(ctx, userID)
	default:
		return []domains.FriendRequest{}, fmt.Errorf("invalid direction '%s': %w", *direction, core_errors.ErrInvalidArg)
	}

	if err != nil {
		return []domains.FriendRequest{}, fmt.Errorf("failed to get friend requests: %w", err)
	}

	return requests, nil
}

func (s *FriendRequestsService) DeleteFriendRequest(ctx context.Context, userID int, requestID int) error {
	if requestID <= 0 {
		return fmt.Errorf("request id must be positive")
	}

	deletedRequest, err := s.friendRequestsRep.Delete(ctx, userID, requestID)
	if err != nil {
		return err
	}

	s.broadcaster.NotifyDeclinedRequest(userID, deletedRequest.FromUser.ID)
	s.broadcaster.NotifyDeclinedRequest(userID, deletedRequest.ToUser.ID)

	return nil
}
