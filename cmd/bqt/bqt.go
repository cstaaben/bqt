/*
	bqt - BigQuery TUI
	Copyright (C) 2024  Corbin Staaben<cstaaben@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/cstaaben/bqt/internal/client"
	"github.com/cstaaben/bqt/internal/config"
	"github.com/cstaaben/bqt/internal/formatter"
	"github.com/cstaaben/bqt/internal/model"
)

var flags = flag.NewFlagSet("bqt", flag.ExitOnError)

func init() {
	flags.StringP("config", "c", os.ExpandEnv("$HOME/.bqt/config"), "Path to the config file.")
	flags.StringP("format", "f", "table", "Format to display query results in.")
	flags.IntP("verbosity", "v", int(slog.LevelInfo), "Verbosity level of logs.")
	flags.String("project-id", "", "GCP project ID to connect to")
	flags.String(
		"credentials-file",
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		"Filepath to the credentials file to use, if different from GOOGLE_APPLICATION_CREDENTIALS.",
	)
	flags.Bool("batch-priority", true, "Run queries with batch priority.")
	flags.DurationP(
		"timeout",
		"t",
		0,
		"Timeout for queries. If 0, queries will run for BigQuery's maximum allowed time.",
	)

	if err := viper.BindPFlags(flags); err != nil {
		panic(fmt.Errorf("binding flags: %w", err))
	}
}

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalln("parsing flags:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.New(viper.GetString("config"))
	if err != nil {
		log.Fatalln("parsing config:", err)
	}

	logFile, err := tea.LogToFile("bqt.log", "[bqt]")
	if err != nil {
		log.Fatalln("error opening log file:", err)
	}
	defer logFile.Close() // nolint:errcheck

	slog.SetDefault(slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.Level(cfg.Verobosity)})))

	slog.DebugContext(ctx, "starting bqt")
	defer slog.DebugContext(ctx, "finished")

	client := client.New()
	defer client.Close() // nolint:errcheck

	var f formatter.Formatter
	switch cfg.Format {
	case "table":
		f = formatter.NewTable()
	case "csv":
		f = formatter.NewCSV()
	case "json":
		f = formatter.NewJSON()
	}

	program := tea.NewProgram(model.New(client, f), tea.WithContext(ctx), tea.WithAltScreen())
	_, err = program.Run()
	if err != nil {
		slog.ErrorContext(ctx, "unexpected error", slog.Any("error", err))
		os.Exit(1)
	}
}
