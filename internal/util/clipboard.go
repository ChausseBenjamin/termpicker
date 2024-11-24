package util

import (
	"fmt"
	"log/slog"

	"github.com/atotto/clipboard"
)

// Copies any object that has the Stringer interface to the clipboard
func Copy(str string) string {
	if err := clipboard.WriteAll(str); err != nil {
		slog.Error("Unable to copy item", "item", str, ErrKey, err)
		return fmt.Sprintf("Copy operation failed: %v", err)
	}
	return fmt.Sprintf("Copied %s to clipboard!", str)
}
