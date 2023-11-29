package types

import (
	"log/slog"
	"strings"
)

type SlogLevel struct {
	Value slog.Level
}

func (l *SlogLevel) Set(value string) error {
	switch strings.ToLower(value) {
	case "error":
		l.Value = slog.LevelError
	case "warn", "warning":
		l.Value = slog.LevelWarn
	case "debug":
		l.Value = slog.LevelDebug
	default:
		l.Value = slog.LevelInfo
	}

	return nil
}
