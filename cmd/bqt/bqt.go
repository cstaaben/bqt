package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cstaaben/bqt/internal/ui"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	verbosity       int
	projectID       string
	credentialsFile string
	batchPriority   bool

	flags = flag.NewFlagSet("bqt", flag.ExitOnError)
)

func init() {
	flags.IntVarP(&verbosity, "verbosity", "v", int(slog.LevelInfo), "Verbosity level of logs.")
	flags.StringVar(&projectID, "project-id", "", "GCP project ID to connect to")
	flags.StringVar(&credentialsFile, "credentials-file", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), "Filepath to the credentials file to use, if different from GOOGLE_APPLICATION_CREDENTIALS.")
	flags.BoolVar(&batchPriority, "batch-priority", true, "Run queries with batch priority.")

	if err := viper.BindPFlags(flags); err != nil {
		panic(fmt.Errorf("binding flags: %w", err))
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.Level(verbosity)}))
	logger.DebugContext(ctx, "starting bqt")
	defer logger.DebugContext(ctx, "finished")

	program := tea.NewProgram(ui.NewModel(ctx, logger))
	_, err := program.Run()
	if err != nil {
		logger.ErrorContext(ctx, "unexpected error", slog.Any("error", err))
		os.Exit(1)
	}
}
