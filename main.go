package main

import (
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/semka95/balance-service/cmd"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	config, err := cmd.NewConfig()
	if err != nil {
		logger.Error("can't decode config", zap.Error(err))
		return
	}

	srv := cmd.NewServer(logger, config)
	srv.RunServer()
}
