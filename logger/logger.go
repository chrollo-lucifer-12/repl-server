package logger

import "log/slog"

type Logger interface {
	Log(level slog.Level, msg string, args ...any)
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}
