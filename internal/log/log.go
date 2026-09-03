// Package log provides centralized logging for the entire application.
package log

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/nightails/leafy/internal/config"
)

type CloseFunc func() error

func InitializeLogger(cfg *config.Config) (*slog.Logger, CloseFunc, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0700); err != nil {
		return nil, nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}
	bufferFile := bufio.NewWriterSize(file, 8192)

	closeLogger := func() error {
		return file.Close()
	}

	return slog.New(slog.NewJSONHandler(bufferFile, &slog.HandlerOptions{Level: slog.LevelInfo})), closeLogger, nil
}
