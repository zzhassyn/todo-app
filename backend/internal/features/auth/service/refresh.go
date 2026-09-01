package auth_service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

type RefreshResult struct {
	User         domain.User
	Token        string
	RefreshToken string
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (RefreshResult, error) {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	userID, err := s.authRepository.GetUserIDFromRefreshToken(ctx, tokenHash)
	if err != nil {
		return RefreshResult{}, fmt.Errorf("validate refresh token: %w", err)
	}

	user, err := s.usersRegistry.GetUser(ctx, userID)
	if err != nil {
		return RefreshResult{}, fmt.Errorf("get user: %w", err)
	}

	// Rotate the refresh token by deleting the old one
	_ = s.authRepository.DeleteRefreshToken(ctx, tokenHash)

	// Issue new access and refresh tokens
	newAccessToken, newRefreshToken, err := s.issueTokens(ctx, user)
	if err != nil {
		return RefreshResult{}, fmt.Errorf("issue tokens: %w", err)
	}

	return RefreshResult{User: user, Token: newAccessToken, RefreshToken: newRefreshToken}, nil
}
