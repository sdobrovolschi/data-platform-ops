package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	log := slog.New(jsonHandler)

	cleaner, err := NewCleaner(log)
	if err != nil {
		log.ErrorContext(ctx, "bootstrap failed", "err", err)
		return
	}

	defer func() {
		if err := cleaner.Close(); err != nil {
			log.ErrorContext(ctx, "shutdown failed", "err", err)
		}
	}()

	if err := cleaner.run(ctx); err != nil {
		log.ErrorContext(ctx, "execution failed", "err", err)
	}
}
