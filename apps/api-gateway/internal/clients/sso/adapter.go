package sso

import (
	"log/slog"

	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/sso/v1"
)

type Adapter struct {
	log *slog.Logger
	api ssov1.AuthClient
}

func New(log *slog.Logger, api ssov1.AuthClient) *Adapter {
	return &Adapter{
		log: log,
		api: api,
	}
}
