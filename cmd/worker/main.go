package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("🚀 System Monitor is running...")

	// يمكن وضع منطق المراقبة أو التنسيق هنا
	logger.Info("📡 Watching all events...")
	select {}
}
