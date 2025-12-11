package tasks

import (
	"log/slog"
)

type Adapter struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Adapter {
	return &Adapter{
		log: log,
	}
}
