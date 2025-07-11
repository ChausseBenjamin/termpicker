package app

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	"github.com/ChausseBenjamin/termpicker/internal/logging"
	"github.com/ChausseBenjamin/termpicker/internal/parse"
	"github.com/ChausseBenjamin/termpicker/internal/preview"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/urfave/cli/v3"
)

//go:embed description.txt
var Desc string

func AppAction(ctx context.Context, cmd *cli.Command) error {
	logfile := logging.Setup(cmd.String(flagLogfile))
	defer logfile.Close()

	slog.Info("Starting Termpicker")

	sw := switcher.New()

	if colorStr := cmd.String(flagColor); colorStr != "" {
		sw.NewNotice(sw.SetColorFromText(colorStr))
	}

	previewStr := cmd.String(flagSampleStr)
	fg := cmd.String(flagSampleFG)
	bg := cmd.String(flagSampleBG)
	for _, ptr := range []*string{&fg, &bg} {
		if *ptr != "" { // <- Hack: keeps original terminal colors if none given
			cs, err := parse.Color(*ptr)
			if err != nil {
				fmt.Println("Failed to parse preview Color")
				os.Exit(1)
			}
			val := colors.Hex(cs)
			ptr = &val
		}
	}
	cfg := preview.Config{
		PreviewStr: previewStr,
		PreviewBg:  bg,
		PreviewFg:  fg,
	}
	sw.UpdatePreview(cfg)

	p := tea.NewProgram(sw,
		tea.WithAltScreen(),
		tea.WithColorProfile(colorprofile.TrueColor),
	)
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
