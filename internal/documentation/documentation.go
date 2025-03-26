/*
 * This package isn't the actual termpicker app.
 * To avoid importing packages which aren't needed at runtime,
 * some auto-generation functionnalities is offloaded to here so
 * it can be done with access to the rest of the code-base but
 * without bloating the final binary. For example,
 * generating bash+zsh auto-completion scripts isn't needed in
 * the final binary if those script are generated before hand.
 * Same gose for manpages. This file is meant to be run automatically
 * to easily package new releases. Same goes for manpages which is the
 * only feature currently in here.
 */
package main

import (
	_ "embed"
	"log/slog"
	"os"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/app"
	docs "github.com/urfave/cli-docs/v3"
)

func main() {
	a := app.Command()

	docs.MarkdownDocTemplate = strings.Join(
		[]string{
			docs.MarkdownDocTemplate,
			app.KeybindingDocs,
		},
		"\n",
	)

	man, err := docs.ToManWithSection(a, 1)
	if err != nil {
		slog.Error("failed to generate man page",
			slog.Any("error_message", err),
		)
		os.Exit(1)
	}
	os.Stdout.Write([]byte(man))
}
