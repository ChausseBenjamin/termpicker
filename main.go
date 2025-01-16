package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/logging"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
)

// Set by the build system
var version = "compiled"

func AppAction(ctx context.Context, cmd *cli.Command) error {
	logfile := logging.Setup(cmd.String("logfile"))
	defer logfile.Close()

	slog.Info("Starting Termpicker")

	sw := switcher.New()

	if colorStr := cmd.String("color"); colorStr != "" {
		sw.NewNotice(sw.SetColorFromText(colorStr))
	}

	p := tea.NewProgram(sw)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.Command{
		Name:                  "Termpicker",
		Usage:                 "A terminal-based color picker",
		Action:                AppAction,
		Authors:               []any{"Benjamin Chausse <benjamin@chausse.xhz>"},
		Version:               version,
		Flags:                 AppFlags,
		EnableShellCompletion: true,
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("Program crashed", util.ErrKey, err.Error())
		os.Exit(1)
	}
}
