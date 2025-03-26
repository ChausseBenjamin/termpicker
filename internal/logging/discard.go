package logging

import (
	"context"
	"log/slog"
)

// DiscardHandler discards all log output. DiscardHandler.Enabled returns false for all Levels.
type DiscardHandler struct{}

func (d DiscardHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (d DiscardHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (d DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return d
}

func (d DiscardHandler) WithGroup(name string) slog.Handler {
	return d
}
