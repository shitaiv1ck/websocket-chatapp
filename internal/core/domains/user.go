package domains

import (
	"fmt"

	core_errors "github.com/shitaiv1ck/realtime-chat/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Username     string
	Password     string
	PasswordHash string
	IsOnline     bool
}

func (u *User) Validate() error {
	if len([]rune(u.Username)) < 3 || len([]rune(u.Username)) > 100 {
		return fmt.Errorf("len 'username' must be between 3 and 100: %w", core_errors.ErrInvalidArg)
	}

	if len([]rune(u.Password)) < 8 || len([]rune(u.Password)) > 100 {
		return fmt.Errorf("len 'password' must be between 3 and 100: %w", core_errors.ErrInvalidArg)
	}

	return nil
}

func (u *User) GeneratePasswordHash() error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)

	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))

	return err == nil
}

type UserPatch struct {
	Username    Nullable[string]
	OldPassword Nullable[string]
	NewPassword Nullable[string]
}

func (u *UserPatch) Validate() error {
	if u.Username.Set && u.Username.Value == nil {
		return fmt.Errorf("'username' can't be null: %w", core_errors.ErrInvalidArg)
	}

	if u.OldPassword.Set && u.OldPassword.Value == nil {
		return fmt.Errorf("'old_password' can't be null: %w", core_errors.ErrInvalidArg)
	}

	if u.NewPassword.Set && u.NewPassword.Value == nil {
		return fmt.Errorf("'new_password' can't be null: %w", core_errors.ErrInvalidArg)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	temp := *u

	if patch.Username.Set {
		if len([]rune(*patch.Username.Value)) < 3 || len([]rune(*patch.Username.Value)) > 100 {
			return fmt.Errorf("len 'username' must be between 3 and 100: %w", core_errors.ErrInvalidArg)
		}

		temp.Username = *patch.Username.Value
	}

	if patch.NewPassword.Set {
		if len([]rune(*patch.NewPassword.Value)) < 8 || len([]rune(*patch.NewPassword.Value)) > 100 {
			return fmt.Errorf("len 'password' must be between 8 and 100: %w", core_errors.ErrInvalidArg)
		}

		temp.Password = *patch.NewPassword.Value
		if err := temp.GeneratePasswordHash(); err != nil {
			return fmt.Errorf("failed to generate password hash")
		}
	}

	*u = temp

	return nil
}
