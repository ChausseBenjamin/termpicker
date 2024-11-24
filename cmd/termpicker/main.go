package main

import (
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/picker"
	"github.com/ChausseBenjamin/termpicker/internal/slider"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	"github.com/ChausseBenjamin/termpicker/internal/util"
	"github.com/charmbracelet/bubbles/progress"
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
	// RGB {{{
	r := slider.New('R', 255, progress.WithGradient("#660000", "#ff0000"))
	g := slider.New('G', 255, progress.WithGradient("#006600", "#00ff00"))
	b := slider.New('B', 255, progress.WithGradient("#000066", "#0000ff"))
	rgb := picker.New([]slider.Model{r, g, b}, "RGB")
	// }}}
	// CYMK {{{
	c := slider.New('C', 100, progress.WithGradient("#006666", "#00ffff"))
	m := slider.New('M', 100, progress.WithGradient("#660066", "#ff00ff"))
	y := slider.New('Y', 100, progress.WithGradient("#666600", "#ffff00"))
	k := slider.New('K', 100, progress.WithSolidFill("#000000"))
	cmyk := picker.New([]slider.Model{c, m, y, k}, "CMYK")
	// }}}
	// HSL {{{
	h := slider.New('H', 360, progress.WithDefaultGradient())
	s := slider.New('S', 100, progress.WithDefaultGradient())
	l := slider.New('L', 100, progress.WithGradient("#222222", "#ffffff"))
	hsl := picker.New([]slider.Model{h, s, l}, "HSL")
	// }}}
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
