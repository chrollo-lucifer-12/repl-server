package logger

import (
	"context"
	"log/slog"
	"os"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger() Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return &SlogLogger{logger: logger}
}

func (s *SlogLogger) Log(level slog.Level, msg string, args ...any) {
	s.logger.Log(context.Background(), level, msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	s.logger.Error(msg, args...)
}
func (s *SlogLogger) Info(msg string, args ...any) {
	s.logger.Info(msg, args...)
}
