package logwarts

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	Info  func(ctx context.Context, message string, attributes ...slog.Attr)
	Error func(ctx context.Context, message string, attributes ...slog.Attr)
	Warn  func(ctx context.Context, message string, attributes ...slog.Attr)
	Debug func(ctx context.Context, message string, attributes ...slog.Attr)
}

func GetLogger(defaultKeys []string, customLogger *slog.Logger) Logger {
	logger := customLogger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	getDefaultAttributes := func(ctx context.Context) []slog.Attr {
		var defAttrs []slog.Attr
		for _, key := range defaultKeys {
			if value := ctx.Value(key); value != nil {
				defAttrs = append(defAttrs, slog.Any(key, value))
			}
		}
		return defAttrs
	}

	logMessage := func(ctx context.Context, level slog.Level, message string, attributes ...slog.Attr) {
		attr := append(getDefaultAttributes(ctx), attributes...)
		logger.LogAttrs(ctx, level, message, attr...)
	}

	return Logger{
		Info: func(ctx context.Context, message string, attributes ...slog.Attr) {
			logMessage(ctx, slog.LevelInfo, message, attributes...)
		},
		Warn: func(ctx context.Context, message string, attributes ...slog.Attr) {
			logMessage(ctx, slog.LevelWarn, message, attributes...)
		},
		Debug: func(ctx context.Context, message string, attributes ...slog.Attr) {
			logMessage(ctx, slog.LevelDebug, message, attributes...)
		},
		Error: func(ctx context.Context, message string, attributes ...slog.Attr) {
			logMessage(ctx, slog.LevelError, message, attributes...)
		},
	}
}
