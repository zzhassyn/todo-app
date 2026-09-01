package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, tokenHash string, userID int, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.refresh_tokens (token_hash, user_id, expires_at)
	VALUES ($1, $2, $3)
	`
	_, err := r.pool.Exec(ctx, query, tokenHash, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetUserIDFromRefreshToken(ctx context.Context, tokenHash string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT user_id, expires_at
	FROM todoapp.refresh_tokens
	WHERE token_hash = $1
	`
	var userID int
	var expiresAt time.Time
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("refresh token not found: %w", core_errors.ErrNotFound)
		}
		return 0, fmt.Errorf("scan refresh token: %w", err)
	}

	if time.Now().UTC().After(expiresAt) {
		return 0, fmt.Errorf("refresh token expired: %w", core_errors.ErrUnauthorized)
	}

	return userID, nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.refresh_tokens
	WHERE token_hash = $1
	`
	_, err := r.pool.Exec(ctx, query, tokenHash)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

// DeleteExpiredRefreshTokens removes all refresh tokens whose expiry has
// passed. Returns the number of rows deleted so the caller can log it.
func (r *AuthRepository) DeleteExpiredRefreshTokens(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM todoapp.refresh_tokens WHERE expires_at < NOW()`
	tag, err := r.pool.Exec(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("delete expired refresh tokens: %w", err)
	}
	return tag.RowsAffected(), nil
}

