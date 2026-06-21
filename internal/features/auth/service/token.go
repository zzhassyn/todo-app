package auth_service

import (
	"fmt"

	core_auth "github.com/zzhassyn/todo-app/internal/core/auth"
	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (s *AuthService) issueToken(user domain.User) (string, error) {
	return core_auth.IssueToken(s.config.JWTSecret, user.ID, user.Email, s.config.TokenTTL)
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
