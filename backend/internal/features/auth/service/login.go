package auth_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	User         domain.User
	Token        string
	RefreshToken string
}

// invalidCredentialsMessage is intentionally identical whether the email
// does not exist or the password is wrong, so the API does not leak which
// part of the credentials was incorrect (an account-enumeration safeguard).
const invalidCredentialsMessage = "invalid email or password"

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	plaintextPassword string,
) (LoginResult, error) {
	user, err := s.usersRegistry.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return LoginResult{}, fmt.Errorf("%s: %w", invalidCredentialsMessage, core_errors.ErrUnauthorized)
		}

		return LoginResult{}, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plaintextPassword)); err != nil {
		return LoginResult{}, fmt.Errorf("%s: %w", invalidCredentialsMessage, core_errors.ErrUnauthorized)
	}

	accessToken, refreshToken, err := s.issueTokens(ctx, user)
	if err != nil {
		return LoginResult{}, fmt.Errorf("issue tokens: %w", err)
	}

	return LoginResult{User: user, Token: accessToken, RefreshToken: refreshToken}, nil
}
