package main

import (
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/logging"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func AppAction(ctx *cli.Context) error {
	logfile := logging.Setup(ctx.String("logfile"))
	defer logfile.Close()

	slog.Info("Starting Termpicker")

	sw := switcher.New()

	if colorStr := ctx.String("color"); colorStr != "" {
		sw.NewNotice(sw.SetColorFromText(colorStr))
	}

	p := tea.NewProgram(sw)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:   "Termpicker",
		Usage:  "A terminal-based color picker",
		Action: AppAction,
		Flags:  AppFlags,
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("Program crashed", util.ErrKey, err.Error())
		os.Exit(1)
	}
}
