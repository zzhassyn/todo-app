package auth_service

import (
	"context"
	"time"

	"github.com/zzhassyn/todo-app/internal/core/domain"
)

// UsersRegistry is the subset of the users feature's service that the auth
// feature depends on. Auth does not import the users repository directly;
// it depends on this interface, satisfied by users_service.UsersService.
// This keeps the cross-feature dependency inverted and explicit, the same
// pattern used by tasks_service.UsersChecker.
type UsersRegistry interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, tokenHash string, userID int, expiresAt time.Time) error
	GetUserIDFromRefreshToken(ctx context.Context, tokenHash string) (int, error)
	DeleteRefreshToken(ctx context.Context, tokenHash string) error
}

type Config struct {
	JWTSecret       string
	TokenTTL        time.Duration
	RefreshTokenTTL time.Duration
}

type AuthService struct {
	authRepository AuthRepository
	usersRegistry  UsersRegistry
	config         Config
}

func NewAuthService(authRepository AuthRepository, usersRegistry UsersRegistry, config Config) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		usersRegistry:  usersRegistry,
		config:         config,
	}
}
