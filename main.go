package main

import (
	"go.uber.org/zap"
	"joborder/config"
	"joborder/internal/server"
	"joborder/pkg/postgres"
)

func main() {

	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(logger)

	getConfig, err := config.GetConfig()
	if err != nil {
		return
	}

	conn, err := postgres.InitDatabaseConn(getConfig)
	if err != nil {
		return
	}

	s := server.NewServer(conn,logger)

	serverErr := s.Run()
	if serverErr != nil {
		logger.Sugar().Error(serverErr)
		return
	}

}
