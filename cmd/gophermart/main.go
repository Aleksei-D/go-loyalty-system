package main

import (
	"github.com/Aleksei-D/go-loyalty-system/internal/app"
	"github.com/Aleksei-D/go-loyalty-system/internal/config"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"go.uber.org/zap"
)

func main() {
	err := logger.Initialize("INFO")
	if err != nil {
		logger.Log.Fatal("cannot initialize zap", zap.Error(err))
	}

	configServe, err := config.NewServerConfig()
	if err != nil {
		logger.Log.Fatal("cannot initialize config", zap.Error(err))
	}

	gopherMartApp, err := app.NewApp(configServe)
	if err != nil {
		logger.Log.Fatal("cannot initialize gophermart", zap.Error(err))
	}

	err = gopherMartApp.Run()
	if err != nil {
		logger.Log.Fatal("cannot start gophermart", zap.Error(err))
	}
}
