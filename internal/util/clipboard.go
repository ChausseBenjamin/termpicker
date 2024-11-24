package util

import (
	"log/slog"

	"github.com/atotto/clipboard"
)

// Copies any object that has the Stringer interface to the clipboard
func Copy(str string) {
	if err := clipboard.WriteAll(str); err != nil {
		slog.Error("Unable to copy item", "item", str, ErrKey, err)
	}
}
