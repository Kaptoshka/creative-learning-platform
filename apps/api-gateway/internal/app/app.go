package app

import (
	"log/slog"
)

type App struct {
	log *slog.Logger
}

func New(log *slog.Logger) *App {
	return &App{
		log: log,
	}
}
