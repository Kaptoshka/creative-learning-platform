package service

import (
	"encoding/base64"
	"strconv"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
)

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return domain.DefaultPageSizeLimit
	}
	if limit > domain.MaxPageSizeLimit {
		return domain.MaxPageSizeLimit
	}
	return limit
}

func EncodePageToken(offset, returned, limit int) string {
	if returned < limit {
		return ""
	}
	return base64.StdEncoding.EncodeToString(
		[]byte(strconv.Itoa(offset + returned)),
	)
}

func DecodePageToken(token string) (int, error) {
	if token == "" {
		return 0, nil
	}
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, domain.ErrInvalidPageToken
	}
	offset, err := strconv.Atoi(string(b))
	if err != nil || offset < 0 {
		return 0, domain.ErrInvalidPageToken
	}
	return offset, nil
}
