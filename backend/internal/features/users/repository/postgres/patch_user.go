package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zzhassyn/todo-app/internal/core/domain"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.users
		SET
			full_name=$1,
			phone_number=$2,
			email=$3,
			password_hash=$4,
			version=version+1
		WHERE id=$5 AND version=$6
		RETURNING id, version, full_name, phone_number, email, password_hash;
	`

	row := r.pool.QueryRow(ctx, query,
		user.FullName,
		user.PhoneNumber,
		user.Email,
		user.PasswordHash,
		id,
		user.Version,
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
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user with id='%d' concurrently accessed: %w", id, core_errors.ErrConflict)
		}

		if errors.Is(err, core_postgres_pool.ErrUniqueViolation) {
			return domain.User{}, fmt.Errorf("user with email='%s': %w", user.Email, core_errors.ErrConflict)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	return userDomainFromModel(userModel), nil
}
