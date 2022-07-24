package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/renanqts/xdpdropper/pkg/api"
	"github.com/renanqts/xdpdropper/pkg/config"
	"github.com/renanqts/xdpdropper/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logConfig := logger.NewDefaultConfig()
	logConfig.Level = config.LogLevel
	logger.Init(logConfig)

	api, err := api.New(config)
	if err != nil {
		logger.Log.Fatal("Failed to initialize api", zap.Error(err))
	}

	if err = api.Start(); err != nil {
		api.Close()
		logger.Log.Fatal("Failed to start api", zap.Error(err))
	}

	// wait for system signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Log.Info("XDPDropper stopping")
	api.Close()
}
