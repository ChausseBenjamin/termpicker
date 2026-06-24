package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ChausseBenjamin/termpicker/internal/util"
)

var globalLvl slog.LevelVar

func Setup(filepath string, level string) *os.File {
	var file *os.File = nil

	if filepath != "" {
		logFile, err := os.Create(filepath)
		if err != nil {
			slog.Error("Failed to create log file", util.ErrKey, err.Error())
			os.Exit(1)
		}

		handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: &globalLvl})
		slog.SetDefault(slog.New(handler))
		file = logFile
	} else {
		// Since app is a TUI, logging to stdout/stderr would break the UI
		// So we disable it by default
		handler := DiscardHandler{}
		slog.SetDefault(slog.New(handler))
	}
	globalLvl.Set(parseLvl(level))
	return file
}

func parseLvl(input string) slog.Level {
	if input == "" {
		return slog.LevelInfo
	}
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
