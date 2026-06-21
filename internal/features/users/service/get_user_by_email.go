package users_service

import (
	"context"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

func (s *UsersService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := s.usersRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by email from repository: %w", err)
	}

	return user, nil
}
