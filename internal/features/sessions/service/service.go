package sessions_service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/shitaiv1ck/realtime-chat/internal/core/domains"
	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
)

type SessionsService struct {
	sessionsRep SessionsRepository
	usersRep    UsersRepository
}

type SessionsRepository interface {
	Save(ctx context.Context, session domains.Session) (domains.Session, error)
	FindByToken(ctx context.Context, token string) (domains.Session, error)
	Delete(ctx context.Context, token string) error
}

type UsersRepository interface {
	FindByUsername(ctx context.Context, username string) (domains.User, error)
}

func NewService(sessionsRep SessionsRepository, usersRep UsersRepository) *SessionsService {
	return &SessionsService{
		sessionsRep: sessionsRep,
		usersRep:    usersRep,
	}
}

func (s *SessionsService) Authentication(ctx context.Context, user domains.User) (int, error) {
	if err := user.Validate(); err != nil {
		return -1, fmt.Errorf("failed to validate user: %w", err)
	}

	foundUser, err := s.usersRep.FindByUsername(ctx, user.Username)
	if err != nil {
		return -1, core_errors.ErrUnauthenticate
	}

	if !foundUser.ComparePassword(user.Password) {
		return -1, core_errors.ErrUnauthenticate
	}

	return foundUser.ID, nil
}

func (s *SessionsService) CreateSession(ctx context.Context, user domains.User) (domains.Session, error) {
	userID, err := s.Authentication(ctx, user)
	if err != nil {
		return domains.Session{}, fmt.Errorf("failed to authenticate: %w", err)
	}

	sessionToken, err := generateToken(32)
	if err != nil {
		return domains.Session{}, fmt.Errorf("failed to generate cookie token: %w", err)
	}

	csrfToken, err := generateToken(32)
	if err != nil {
		return domains.Session{}, fmt.Errorf("failed to generate cookie token: %w", err)
	}

	session := domains.Session{
		SessionToken: sessionToken,
		CSRFToken:    csrfToken,
		UserID:       userID,
		ExpiredAt:    time.Now().Add(24 * time.Hour),
	}

	if err := session.Validate(); err != nil {
		return domains.Session{}, fmt.Errorf("failed to validate session: %w", err)
	}

	createdSession, err := s.sessionsRep.Save(ctx, session)
	if err != nil {
		return domains.Session{}, err
	}

	return createdSession, nil
}

func (s *SessionsService) GetSession(ctx context.Context, sessionToken string) (domains.Session, error) {
	session, err := s.sessionsRep.FindByToken(ctx, sessionToken)
	if err != nil {
		return domains.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	if err := session.Validate(); err != nil {
		return domains.Session{}, fmt.Errorf("failed to validate session: %w", err)
	}

	return session, nil
}

func (s *SessionsService) DeleteSession(ctx context.Context, token string) error {
	if err := s.sessionsRep.Delete(ctx, token); err != nil {
		return err
	}

	return nil
}

func generateToken(len int) (string, error) {
	bytes := make([]byte, len)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)

	return token, nil
}
