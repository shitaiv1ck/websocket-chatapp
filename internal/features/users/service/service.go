package users_service

import (
	"context"
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type UsersService struct {
	rep         UsersRepository
	broadcaster UsersWSTransport
}

type UsersRepository interface {
	Save(ctx context.Context, user domains.User) (domains.User, error)
	FindAll(ctx context.Context, limit *int, offset *int) ([]domains.User, error)
	FindByID(ctx context.Context, id int) (domains.User, error)
	Update(ctx context.Context, user domains.User) (domains.User, error)
}

type UsersWSTransport interface {
	BroadcastEvent(event string, user domains.User)
}

func NewService(rep UsersRepository, broadcaster UsersWSTransport) *UsersService {
	return &UsersService{
		rep:         rep,
		broadcaster: broadcaster,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, user domains.User) (domains.User, error) {
	if err := user.Validate(); err != nil {
		return domains.User{}, fmt.Errorf("filed to validate user: %w", err)
	}

	if err := user.GeneratePasswordHash(); err != nil {
		return domains.User{}, fmt.Errorf("failed to generate password hash: %w", err)
	}

	createdUser, err := s.rep.Save(ctx, user)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	s.broadcaster.BroadcastEvent("user.created", createdUser)

	return createdUser, nil
}

func (s *UsersService) GetUsers(ctx context.Context, limit *int, offset *int) ([]domains.User, error) {
	if limit != nil && *limit < 0 {
		return []domains.User{}, fmt.Errorf("limit must be non negative: %w", core_errors.ErrInvalidArg)
	}

	if offset != nil && *offset < 0 {
		return []domains.User{}, fmt.Errorf("offset must be non negative: %w", core_errors.ErrInvalidArg)
	}

	users, err := s.rep.FindAll(ctx, limit, offset)
	if err != nil {
		return []domains.User{}, fmt.Errorf("failed to get users from repository: %w", err)
	}

	return users, err
}

func (s *UsersService) GetUser(ctx context.Context, userID int) (domains.User, error) {
	user, err := s.rep.FindByID(ctx, userID)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to get user from repository: %w", err)
	}

	return user, nil
}

func (s *UsersService) PatchUser(ctx context.Context, userID int, patch domains.UserPatch) (domains.User, error) {
	user, err := s.rep.FindByID(ctx, userID)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to get user from repository: %w", err)
	}

	if err := patch.Validate(); err != nil {
		return domains.User{}, fmt.Errorf("failed to validate patch: %w", err)
	}

	if patch.OldPassword.Set {
		if !user.ComparePassword(*patch.OldPassword.Value) {
			return domains.User{}, fmt.Errorf("invalid password: %w", core_errors.ErrInvalidArg)
		}
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domains.User{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	patchedUser, err := s.rep.Update(ctx, user)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to patch user: %w", err)
	}

	if patch.Username.Set {
		s.broadcaster.BroadcastEvent("user.change_username", patchedUser)
	}

	return patchedUser, nil
}
