package main

import (
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/picker"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

func AppAction(ctx *cli.Context) error {
	logFile, err := os.Create(ctx.String(flagLogfile))
	if err != nil {
		slog.Error("Failed to create log file", util.ErrKey, err.Error())
		os.Exit(1)
	}
	defer logFile.Close()

	handler := slog.NewJSONHandler(logFile, nil)
	slog.SetDefault(slog.New(handler))

	slog.Info("Starting Termpicker")
	rgb := picker.RGB()
	cmyk := picker.CMYK()
	hsl := picker.HSL()
	sw := switcher.New([]picker.Model{*rgb, *cmyk, *hsl})
	p := tea.NewProgram(sw)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:   "TermPicker",
		Usage:  "A terminal-based color picker",
		Action: AppAction,
		Flags:  AppFlags,
	}
	if err := app.Run(os.Args); err != nil {
		slog.Error("Program crashed", util.ErrKey, err.Error())
		os.Exit(1)
	}
}
