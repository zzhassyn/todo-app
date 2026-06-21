package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/zzhassyn/todo-app/internal/core/errors"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("param='%s' by key='%s' not a valid integer: %v: %w",
			value, key, err, core_errors.ErrInvalidArgument)
	}

	return &val, nil
}

func GetBoolQueryParam(r *http.Request, key string) (*bool, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, nil
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("param='%s' by key='%s' not a valid boolean: %v: %w",
			value, key, err, core_errors.ErrInvalidArgument)
	}

	return &val, nil
}
