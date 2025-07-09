package util

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
)

func Copy(str string) string {
	osc52 := fmt.Sprintf("\033]52;c;%s\a", base64.StdEncoding.EncodeToString([]byte(str)))
	_, err := os.Stdout.WriteString(osc52)
	if err != nil {
		slog.Error("OSC52 write failed", "item", str, ErrKey, err)
		return fmt.Sprintf("Failed to copy '%s': OSC52 failed", str)
	}
	return fmt.Sprintf("Copied '%s' using OSC52", str)
}
