package app

import (
	"context"
	_ "embed"
	"log/slog"

	"github.com/ChausseBenjamin/termpicker/internal/logging"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"
)

//go:embed description.txt
var Desc string

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

func Command(version string) *cli.Command {
	cmd := &cli.Command{
		Name:                  "termpicker",
		Usage:                 "A terminal-based color picker",
		Action:                AppAction,
		ArgsUsage:             "",
		Description:           Desc,
		Authors:               []any{"Benjamin Chausse <benjamin@chausse.xyz>"},
		Version:               version,
		Flags:                 AppFlags,
		EnableShellCompletion: true,
	}

	return cmd
}
