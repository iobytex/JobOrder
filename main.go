package main

import (
	"go.uber.org/zap"
	"joborder/config"
	"joborder/internal/model"
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

	pgsqlConn, err := postgres.InitDatabaseConn(getConfig)
	if err != nil {
		return
	}

	migrateErr := pgsqlConn.AutoMigrate(&model.Category{},&model.User{},&model.Product{}, &model.Order{},&model.OrderItem{},&model.Stock{})
	if migrateErr != nil {
		panic(migrateErr)
	}

	s := server.NewServer(pgsqlConn,logger)

	serverErr := s.Run()
	if serverErr != nil {
		logger.Sugar().Error(serverErr)
		return
	}

}
