package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

// GetUUIDQueryParam supports one extra sentinel value beyond a literal
// UUID: "none", which callers can use to mean "explicitly filter for no
// value" (e.g. tasks with no folder) as opposed to "filter not provided".
// Since both "absent" and "explicitly none" need to be distinguishable
// from "a real UUID was given", this returns three states via two return
// values: (nil, false, nil) = not provided; (nil, true, nil) = "none"
// requested; (&id, true, nil) = a real UUID was given.
func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, bool, error) {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil, false, nil
	}

	if value == "none" {
		return nil, true, nil
	}

	val, err := uuid.Parse(value)
	if err != nil {
		return nil, false, fmt.Errorf("param='%s' by key='%s' not a valid UUID: %v: %w",
			value, key, err, core_errors.ErrInvalidArgument)
	}

	return &val, true, nil
}
