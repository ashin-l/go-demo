package common

import (
	"os"
	"os/signal"

	"github.com/ashin-l/go-demo/pkg/logger"
)

func WaitSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	s := <-signals
	logger.Logger().Info("signal:", s.String())
}
