package util

import (
	"log/slog"

	"golang.design/x/clipboard"
)

// Copies any object that has the Stringer interface to the clipboard
func Copy(str string) {
	// Initialize the clipboard
	if err := clipboard.Init(); err != nil {
		slog.Error("failed to initialize clipboard", "error", err)
		return
	}

	clipboard.Write(clipboard.FmtText, []byte(str))
}
