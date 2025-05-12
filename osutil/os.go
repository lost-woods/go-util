package osutil

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var (
	zapLogger, _ = zap.NewProduction()
	log          = zapLogger.Sugar()
)

func GetLogger() *zap.SugaredLogger {
	return log
}

func WaitForCtrlC() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
