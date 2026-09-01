package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number, email, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id, version, full_name, phone_number, email, password_hash;
	`

	row := r.pool.QueryRow(ctx, query,
		user.FullName,
		user.PhoneNumber,
		user.Email,
		user.PasswordHash,
	)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
		&userModel.Email,
		&userModel.PasswordHash,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return domain.User{}, fmt.Errorf("user with email='%s': %w", user.Email, core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return userDomainFromModel(userModel), nil
}
