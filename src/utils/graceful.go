package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown() {
	// Catch signal so we can shutdown gracefully
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	// Wait for a signal
	<-sigCh
}
