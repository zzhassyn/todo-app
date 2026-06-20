package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

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
