package sessions_service

import (
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
	Save(session domains.Session) (domains.Session, error)
	FindByToken(token string) (domains.Session, error)
	Delete(token string) error
}

type UsersRepository interface {
	FindByUsername(username string) (domains.User, error)
}

func NewService(sessionsRep SessionsRepository, usersRep UsersRepository) *SessionsService {
	return &SessionsService{
		sessionsRep: sessionsRep,
		usersRep:    usersRep,
	}
}

func (s *SessionsService) Authentication(user domains.User) (int, error) {
	if err := user.Validate(); err != nil {
		return -1, fmt.Errorf("failed to validate user: %w", err)
	}

	foundUser, err := s.usersRep.FindByUsername(user.Username)
	if err != nil {
		return -1, core_errors.ErrUnauthenticate
	}

	if !foundUser.ComparePassword(user.Password) {
		return -1, core_errors.ErrUnauthenticate
	}

	return foundUser.ID, nil
}

func (s *SessionsService) CreateSession(userID int) (domains.Session, error) {
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
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	createdSession, err := s.sessionsRep.Save(session)
	if err != nil {
		return domains.Session{}, err
	}

	return createdSession, nil
}

func (s *SessionsService) GetSession(sessionToken string) (domains.Session, error) {
	session, err := s.sessionsRep.FindByToken(sessionToken)
	if err != nil {
		return domains.Session{}, fmt.Errorf("failed to get session: %w", err)
	}

	if !session.ExpiresAt.After(time.Now()) {
		return domains.Session{}, fmt.Errorf("expired cookie: %w", core_errors.ErrCoockie)
	}

	return session, nil
}

func (s *SessionsService) DeleteSession(token string) error {
	if err := s.sessionsRep.Delete(token); err != nil {
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
