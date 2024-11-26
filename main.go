package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/ChausseBenjamin/termpicker/internal/logging"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

var ( // Set by the build system
	version = "compiled"
	date    = ""
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
	compileDate, _ := time.Parse(time.RFC3339, date)
	app := &cli.App{
		Name:   "Termpicker",
		Usage:  "A terminal-based color picker",
		Action: AppAction,
		Authors: []*cli.Author{
			{Name: "Benjamin Chausse", Email: "benjamin@chausse.xyz"},
		},
		Version:  version,
		Flags:    AppFlags,
		Compiled: compileDate,
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("Program crashed", util.ErrKey, err.Error())
		os.Exit(1)
	}
}
