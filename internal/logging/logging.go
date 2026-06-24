package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/util"
)

func Setup(filepath string, level string) *os.File {
	if filepath != "" {
		logFile, err := os.Create(filepath)
		if err != nil {
			slog.Error("Failed to create log file", util.ErrKey, err.Error())
			os.Exit(1)
		}

		handler := slog.NewJSONHandler(logFile, nil)
		slog.SetDefault(slog.New(handler))
		if level != "" {
			slog.SetLogLoggerLevel(parseLvl(level))
		}

		return logFile
	} else {
		// Since app is a TUI, logging to stdout/stderr would break the UI
		// So we disable it by default
		handler := DiscardHandler{}
		slog.SetDefault(slog.New(handler))
		return nil
	}
}

func parseLvl(input string) slog.Level {
	input = strings.ToLower(input)
	for _, v := range []struct {
		prefix string
		level  slog.Level
	}{
		{"deb", slog.LevelDebug},
		{"inf", slog.LevelInfo},
		{"war", slog.LevelWarn},
		{"err", slog.LevelError},
	} {
		if strings.HasPrefix(input, v.prefix) {
			return v.level
		}
	}
	slog.Warn("Failed to parse provided log-level value",
		"level", input,
	)
	return slog.LevelInfo
}
