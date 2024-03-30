package ui

import (
	"log/slog"

	"github.com/cstaaben/bqt/internal/client"
)

type Model struct {
	client *client.Client
	logger *slog.Logger

	QueryResults []QueryResult
}

type QueryResult struct{}

func NewModel(client *client.Client, logger slog.Logger) *Model {
	return &Model{
		client: client,
		logger: logger.With("name", "ui"),
	}
}
