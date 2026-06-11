package friendrequests_service

import (
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type FriendRequestsService struct {
	friendRequestsRep FriendRequestsRepository
	usersRep          UsersRepository
	broadcaster       FriendRequestsWSTransport
}

type FriendRequestsRepository interface {
	Save(request domains.FriendRequest) (domains.FriendRequest, error)
	FindByFromIDAndToID(fromID int, toID int) (domains.FriendRequest, error)
	FindByToID(toID int) ([]domains.FriendRequest, error)
	FindByFromID(fromID int) ([]domains.FriendRequest, error)
	Delete(userID int, requestID int) error
}

type UsersRepository interface {
	FindByID(id int) (domains.User, error)
}

type FriendRequestsWSTransport interface {
	NotifyClientEvent(userID int, event string, request domains.FriendRequest)
	NotifyDeclinedRequest(userID int, requestID int)
}

func NewService(friendRequestsRep FriendRequestsRepository, usersRep UsersRepository, broadcaster FriendRequestsWSTransport) *FriendRequestsService {
	return &FriendRequestsService{
		friendRequestsRep: friendRequestsRep,
		usersRep:          usersRep,
		broadcaster:       broadcaster,
	}
}

func (s *FriendRequestsService) CreateFriendRequest(request domains.FriendRequest) (domains.FriendRequest, error) {
	if request.FromUser.ID <= 0 || request.ToUser.ID <= 0 {
		return domains.FriendRequest{}, fmt.Errorf("user id must be positive: %w", core_errors.ErrInvalidArg)
	}

	if request.FromUser.ID == request.ToUser.ID {
		return domains.FriendRequest{}, fmt.Errorf("can't send friend request to yourself: %w", core_errors.ErrInvalidArg)
	}

	fromUser, err := s.usersRep.FindByID(request.FromUser.ID)
	if err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to get user with id=%v: %w", request.FromUser.ID, err)
	}
	toUser, err := s.usersRep.FindByID(request.ToUser.ID)
	if err != nil {
		return domains.FriendRequest{}, fmt.Errorf("failed to get user with id=%v: %w", request.ToUser.ID, err)
	}

	if s.hasIncomingFromTargetUser(request) {
		return domains.FriendRequest{}, fmt.Errorf("user with id=%v already sent friend request: %w", request.ToUser.ID, core_errors.ErrConflict)
	}

	createdFriendRequest, err := s.friendRequestsRep.Save(request)
	if err != nil {
		return domains.FriendRequest{}, err
	}
	createdFriendRequest.FromUser.Username = fromUser.Username
	createdFriendRequest.ToUser.Username = toUser.Username

	s.broadcaster.NotifyClientEvent(request.ToUser.ID, "send_friend_request", createdFriendRequest)

	return createdFriendRequest, nil
}

func (s *FriendRequestsService) hasIncomingFromTargetUser(request domains.FriendRequest) bool {
	_, err := s.friendRequestsRep.FindByFromIDAndToID(request.ToUser.ID, request.FromUser.ID)

	return err == nil
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
