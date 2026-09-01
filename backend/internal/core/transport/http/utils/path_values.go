package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return 0, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.Atoi(pathValue)
	if err != nil {
		return 0, fmt.Errorf("path value='%s' for key='%s' is not a valid integer: %v: %w", pathValue, key, err, core_errors.ErrInvalidArgument)
	}

	return val, nil
}

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	val, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf("path value='%s' for key='%s' is not a valid UUID: %v: %w", pathValue, key, err, core_errors.ErrInvalidArgument)
	}

	return val, nil
}
