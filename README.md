# Logwarts

Logwarts is a simple, context-aware logging package for Go. It provides structured logging with default and custom attributes, using Go's `slog` package.

## Installation

```bash
go get github.com/foxtrot12/logwarts
```

## Features

- Supports `Info`, `Error`, `Warn`, and `Debug` log levels.
- Allows default attributes to be extracted from `context.Context`.
- Works with the default `slog` JSON handler or a custom logger.

## Usage

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/foxtrot12/logwarts"
)

func main() {
	// Define default keys to extract from context
	defaultKeys := []string{"userID", "requestID"}

	// Create a custom logger or use nil for the default logger
	customLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialize Logwarts logger
	logger := logwarts.GetLogger(defaultKeys, customLogger)

	// Create a context with default keys
	ctx := context.WithValue(context.Background(), "userID", "12345")
	ctx = context.WithValue(ctx, "requestID", "abcde")

	// Log messages
	logger.Info(ctx, "User logged in", slog.String("role", "admin"))
	logger.Warn(ctx, "Unusual login attempt detected")
	logger.Error(ctx, "Failed to fetch user data", slog.Int("errorCode", 500))
	logger.Debug(ctx, "Debugging user session", slog.String("sessionID", "xyz123"))
}
```

## License

This project is licensed under the MIT License.
