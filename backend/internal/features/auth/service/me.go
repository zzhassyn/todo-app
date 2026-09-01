package auth_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *AuthService) Me(ctx context.Context, userID int) (domain.User, error) {
	user, err := s.usersRegistry.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
