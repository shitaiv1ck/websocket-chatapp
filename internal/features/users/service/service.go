package users_service

import (
	"fmt"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
)

type UsersService struct {
	rep UsersRepository
}

type UsersRepository interface {
	Save(user domains.User) (domains.User, error)
	FindAll() ([]domains.User, error)
	FindByID(id int) (domains.User, error)
	Update(user domains.User) (domains.User, error)
}

func NewService(rep UsersRepository) *UsersService {
	return &UsersService{
		rep: rep,
	}
}

func (s *UsersService) CreateUser(user domains.User) (domains.User, error) {
	if err := user.Validate(); err != nil {
		return domains.User{}, fmt.Errorf("filed to validate user: %w", err)
	}

	if err := user.GeneratePasswordHash(); err != nil {
		return domains.User{}, fmt.Errorf("failed to generate password hash: %w", err)
	}

	createdUser, err := s.rep.Save(user)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (s *UsersService) GetUsers() ([]domains.User, error) {
	users, err := s.rep.FindAll()
	if err != nil {
		return []domains.User{}, fmt.Errorf("failed to get users from repository: %w", err)
	}

	return users, err
}

func (s *UsersService) PatchUser(id int, patch domains.UserPatch) (domains.User, error) {
	user, err := s.rep.FindByID(id)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to get user from repository: %w", err)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domains.User{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	patchedUser, err := s.rep.Update(user)
	if err != nil {
		return domains.User{}, fmt.Errorf("failed to patch user: %w", err)
	}

	return patchedUser, err
}
