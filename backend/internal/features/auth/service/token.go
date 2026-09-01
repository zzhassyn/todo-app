package auth_service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (s *AuthService) issueTokens(ctx context.Context, user domain.User) (accessToken string, refreshToken string, err error) {
	accessToken, err = core_auth.IssueToken(s.config.JWTSecret, user.ID, user.Email, s.config.TokenTTL)
	if err != nil {
		return "", "", fmt.Errorf("issue token: %w", err)
	}

	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	refreshToken = base64.RawURLEncoding.EncodeToString(rawToken)

	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	expiresAt := time.Now().UTC().Add(s.config.RefreshTokenTTL)
	if err := s.authRepository.CreateRefreshToken(ctx, tokenHash, user.ID, expiresAt); err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func validatePlaintextPassword(password string) error {
	length := len([]rune(password))
	if length < 8 || length > 72 {
		// bcrypt silently truncates inputs over 72 bytes, so an upper bound
		// is enforced here to avoid surprising behavior rather than as a
		// security measure.
		return fmt.Errorf("password must be between 8 and 72 characters long: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
