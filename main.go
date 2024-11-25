package main

import (
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/picker"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	"github.com/ChausseBenjamin/termpicker/internal/userinput"
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

	sw := switcher.New([]picker.Model{
		*picker.RGB(),
		*picker.CMYK(),
		*picker.HSL(),
	})

	if colorStr := ctx.String("color"); colorStr != "" {
		color, err := userinput.ParseColor(colorStr)
		if err != nil {
			slog.Error("Failed to parse color", util.ErrKey, err)
			sw.NewNotice(err.Error())
		} else {
			pc := color.ToPrecise()
			switch color.(type) {
			case colors.RGB:
				sw.UpdatePicker(0, pc)
				sw.SetActive(0)
			case colors.CMYK:
				sw.UpdatePicker(1, pc)
				sw.SetActive(1)
			case colors.HSL:
				sw.UpdatePicker(2, pc)
				sw.SetActive(2)
			}
			sw.NewNotice("Color set to " + colorStr)
		}
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
