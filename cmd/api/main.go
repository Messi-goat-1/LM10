package main

import (
	"LM-Gate/internal/logic"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("🚀 API Server is starting...")

	// تشغيل السيرفر مباشرة
	logic.RunAPIServer()
}
