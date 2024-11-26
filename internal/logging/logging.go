package logging

import (
	"log/slog"
	"os"

	"github.com/ChausseBenjamin/termpicker/internal/util"
)

type logSink struct{}

func (l logSink) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func Setup(filepath string) *os.File {
	if filepath != "" {
		logFile, err := os.Create(filepath)
		if err != nil {
			slog.Error("Failed to create log file", util.ErrKey, err.Error())
			os.Exit(1)
		}

		handler := slog.NewJSONHandler(logFile, nil)
		slog.SetDefault(slog.New(handler))

		return logFile
	} else {
		// Since app is a TUI, logging to stdout/stderr would break the UI
		// So we disable it by default
		handler := slog.NewJSONHandler(logSink{}, nil)
		slog.SetDefault(slog.New(handler))
		return nil
	}
}
