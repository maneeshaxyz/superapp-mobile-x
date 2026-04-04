package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultLimit = 20
	MaxLimit     = 50
)

// ParsePagination parses limit and cursor from the request.
// The cursor is expected to be a base64 encoded string of "timestamp|id".
func ParsePagination(r *http.Request) (limit int, afterID string, afterCreatedAt *time.Time, err error) {
	limitStr := r.URL.Query().Get("limit")
	cursorStr := r.URL.Query().Get("cursor")

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, "", nil, fmt.Errorf("invalid 'limit' parameter: must be an integer")
		}
	}
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	if cursorStr != "" {
		decoded, err := base64.StdEncoding.DecodeString(cursorStr)
		if err != nil {
			return 0, "", nil, fmt.Errorf("invalid 'cursor' parameter: not a valid base64 string")
		}

		parts := strings.Split(string(decoded), "|")
		if len(parts) != 2 {
			return 0, "", nil, fmt.Errorf("invalid 'cursor' parameter: incorrect format")
		}

		ts, err := time.Parse(time.RFC3339, parts[0])
		if err != nil {
			return 0, "", nil, fmt.Errorf("invalid 'cursor' parameter: invalid timestamp")
		}

		afterCreatedAt = &ts
		afterID = parts[1]
	}

	return limit, afterID, afterCreatedAt, nil
}

// EncodeCursor creates a base64 encoded cursor from a timestamp and an ID.
func EncodeCursor(createdAt time.Time, id string) *string {
	payload := fmt.Sprintf("%s|%s", createdAt.Format(time.RFC3339), id)
	encoded := base64.StdEncoding.EncodeToString([]byte(payload))
	return &encoded
}
