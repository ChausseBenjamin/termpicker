package app

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"

	"github.com/ChausseBenjamin/termpicker/internal/logging"
	"github.com/ChausseBenjamin/termpicker/internal/switcher"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

//go:embed keybindings.md
var KeybindingDocs string

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

func Command() *cli.Command {
	cmd := &cli.Command{
		Name:                  "termpicker",
		Usage:                 "A terminal-based color picker",
		Action:                AppAction,
		Authors:               []any{"Benjamin Chausse <benjamin@chausse.xyz>"},
		Version:               version,
		Flags:                 AppFlags,
		EnableShellCompletion: true,
	}

	cli.HelpPrinter = func(w io.Writer, _ string, _ any) {
		docs.MarkdownDocTemplate = fmt.Sprintf("%s\n‎\n\n%s",
			docs.MarkdownDocTemplate,
			KeybindingDocs,
		)

		helpRaw, _ := docs.ToMarkdown(cmd)
		helpCute, _ := glamour.Render(helpRaw, "dark")

		w.Write([]byte(helpCute))
	}

	return cmd
}
