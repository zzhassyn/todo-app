package auth_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

type RegisterResult struct {
	User  domain.User
	Token string
}

func (s *AuthService) Register(
	ctx context.Context,
	fullName string,
	phoneNumber *string,
	email string,
	plaintextPassword string,
) (RegisterResult, error) {
	if err := validatePlaintextPassword(plaintextPassword); err != nil {
		return RegisterResult{}, fmt.Errorf("validate password: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("hash password: %w", err)
	}

	user := domain.NewUserUninitialized(fullName, phoneNumber, email, string(passwordHash))

	user, err = s.usersRegistry.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, core_errors.ErrConflict) {
			return RegisterResult{}, fmt.Errorf("email already registered: %w", err)
		}

		return RegisterResult{}, fmt.Errorf("create user: %w", err)
	}

	token, err := s.issueToken(user)
	if err != nil {
		return RegisterResult{}, fmt.Errorf("issue token: %w", err)
	}

	return RegisterResult{User: user, Token: token}, nil
}
